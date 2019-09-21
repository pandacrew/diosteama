extern crate mysql;
extern crate futures;
extern crate rand; // 0.6.0

extern crate telegram_bot;
extern crate tokio_core;

use std::env;

use futures::Stream;
use tokio_core::reactor::Core;
use telegram_bot::*;
use mysql as my;
use chrono::{DateTime, Utc, TimeZone};
use serde::{Deserialize, Serialize};
use rand::Rng;

#[derive(Debug, Serialize, Deserialize)]
struct Quote {
    recnum: u32,
    quote: String,
    author: String,
    date: DateTime<Utc>,
}

impl Quote {
    fn author_nick(&self) -> &str {
        let i = &self.author.find("!").unwrap();
        let (nick, _) = &self.author.split_at(*i);
        nick
    }
}
impl ::std::fmt::Display for Quote {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        write!(f, "{}", format!("{}\n\n-- Quote {} by {} on {}", &self.quote, &self.recnum, &self.author_nick(), &self.date))
    }
}

fn main() {
    let token = env::var("TELEGRAM_BOT_TOKEN").unwrap();
    let db_url = env::var("DIOSTEAMA_DB_URL").unwrap();
    let pool = my::Pool::new(db_url).unwrap();
    let mut core = Core::new().unwrap();
    println!("{}", rquote(&pool));

    let api = Api::configure(token).build(core.handle()).unwrap();
    let mut rng = rand::thread_rng();
    // Fetch new updates via long poll method
    let future = api.stream().for_each(|update| {

        // If the received update contains a new message...
        if let UpdateKind::Message(message) = update.kind {
            let talk = rng.gen_range(1, 7) == 2;
            if let MessageKind::Text {ref data, ..} = message.kind {
                // Print received text message to stdout.
                println!("<{}>: {:?}", &message.from.first_name, &message);
                let mut iter = data.split_whitespace();
                let cmd = iter.next();
                if cmd == Some("/rquote") ||  cmd == Some("!rquote") {
                    api.spawn(message.chat.text(format!("{}", rquote(&pool))));
                } else if cmd == Some("/quote") || cmd == Some("!quote") {
                    match iter.next() {
                        Some(qid) => match qid.parse::<u32>() {
                            Ok(qid) => api.spawn(message.chat.text(format!("{}", quote_num(&pool, qid)))),
                            Err(_) =>  api.spawn(message.chat.text(format!("{}", quote_like(&pool, &qid.parse::<String>().unwrap())))),
                        }
                        None => api.spawn(message.chat.text(format!("{}", rquote(&pool)))),
                    }
                } else if data.contains("fairy") && talk {
                    api.spawn(message.chat.text("https://ytcropper.com/cropped/ew5cddf953ab0d1"));
                } else if data.contains("dios") && talk {
                    api.spawn(message.chat.text("¿Donde está dios cuando le necesitas, eh? ¿Donde está ese gran maricón misericordioso ahora? Aquí me tienes dios, ¡aquí me tienes cabronazo!\nhttps://www.youtube.com/watch?v=ec2V0tm4JGs"));
                }
            }
        }

        Ok(())
    });
    core.run(future).unwrap();
}

fn quote(pool: &my::Pool, query: &str) -> String {
    let query = format!("SELECT recnum, quote, author, date FROM linux_gey_db {}", query);
    let result: Vec<Quote> =
    pool.prep_exec(query, ())
    .map(|result| {
        result.map(|x| x.unwrap()).map(|row| {
            let (recnum, quote, author, date) = my::from_row(row);
            Quote {
                recnum: recnum,
                quote: quote,
                author: author,
                date: Utc.timestamp(date,0),
            }
        }).collect()
    }).unwrap();
    match result.iter().next() {
        Some(quote) => format!("{}", quote),
        None => format!("Quote no encontrado"),
    }
}

fn quote_like(pool: &my::Pool, needle: &str) -> String {
    let ch = needle.chars().next().unwrap();
    if !ch.is_alphabetic() {
        return format!("{}: Usted necesita la forma A38", needle);
    } else if needle.contains("'") {
        return format!("{}: ¡El puerto sigue estando al lado del mar!", needle);
    }
    let query = format!("WHERE quote LIKE '%{}%' ORDER BY rand() LIMIT 1", str::replace(needle,"*","%"));
    quote(pool, &query)
}

fn quote_num(pool: &my::Pool, qid: u32) -> String {
    let query = format!("WHERE recnum={}", qid);
    quote(pool, &query)
}

fn rquote(pool: &my::Pool) -> String {
    let query = "ORDER BY rand() LIMIT 1";
    quote(pool, &query)
}


// fn export() {
//     serde_json::to_string_pretty(&result[0]).unwrap()
// }
