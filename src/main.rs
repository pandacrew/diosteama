extern crate mysql;
extern crate futures;
extern crate telegram_bot;
extern crate tokio_core;

use std::env;

use futures::Stream;
use tokio_core::reactor::Core;
use telegram_bot::*;
use mysql as my;
use chrono::{DateTime, Utc, TimeZone};
use serde::{Deserialize, Serialize};

struct MiniQuote {
    recnum: u32,
    quote: String,
    author: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Quote {
    recnum: u32,
    quote: String,
    author: String,
    date: DateTime<Utc>,
    deleted: String,
    deleted_by: String,
    deleted_date: DateTime<Utc>,
}
impl MiniQuote {
    fn author_nick(&self) -> &str {
        let i = &self.author.find("!").unwrap();
        let (nick, _) = &self.author.split_at(*i);
        nick
    }
}
impl ::std::fmt::Display for MiniQuote {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        write!(f, "{}", format!("{}\n\n-- Quote {} by {}", &self.quote, &self.recnum, &self.author_nick()))
    }
}

impl ::std::fmt::Display for Quote {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        write!(f, "{}", format!("{}. {}", &self.recnum, &self.quote))
    }
}

fn main() {
    let pool = my::Pool::new("mysql://itorres:7419a533-3dee-4633-9049-c93950378b7a@db1:3306/quotes").unwrap();
    let mut core = Core::new().unwrap();
    println!("{}", rquote(&pool));
    
    let token = env::var("TELEGRAM_BOT_TOKEN").unwrap();
    let api = Api::configure(token).build(core.handle()).unwrap();

    // Fetch new updates via long poll method
    let future = api.stream().for_each(|update| {

        // If the received update contains a new message...
        if let UpdateKind::Message(message) = update.kind {

            if let MessageKind::Text {ref data, ..} = message.kind {
                // Print received text message to stdout.
                println!("<{}>: {:?}", &message.from.first_name, &message);
                let mut iter = data.split_whitespace();
                let cmd = iter.next();
                if cmd == Some("/rquote") {                    
                    api.spawn(message.chat.text(format!("{}", rquote(&pool))));
                } else if cmd == Some("/quote") {
                    match iter.next() {
                        Some(qid) => match qid.parse::<u32>() {
                            Ok(qid) => api.spawn(message.chat.text(format!("{}", quote_num(&pool, qid)))),
                            Err(_) =>  api.spawn(message.chat.text(format!("{}", quote_like(&pool, &qid.parse::<String>().unwrap())))),
                        }
                        None => api.spawn(message.chat.text(format!("{}", rquote(&pool)))),
                    }
                }
            }
        }

        Ok(())
    });
    core.run(future).unwrap();
}

fn quote(pool: &my::Pool, query: &str) -> String {
    let query = format!("SELECT recnum, quote, author FROM linux_gey_db {}", query);    
    let result: Vec<MiniQuote> =
    pool.prep_exec(query, ())
    .map(|result| {
        result.map(|x| x.unwrap()).map(|row| {
            let (recnum, quote, author) = my::from_row(row);
            MiniQuote {
                recnum: recnum,
                quote: quote,
                author: author,
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
        return format!("Usted necesita la forma A38");
    }
    let query = format!("WHERE quote LIKE '%{}%' ORDER BY rand() LIMIT 1", needle);
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
