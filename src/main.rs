extern crate mysql;
extern crate futures;
extern crate rand; // 0.6.0
extern crate telebot;
extern crate tokio_core;

use std::env;

use futures::{stream::Stream};
use telebot::Bot;
use mysql as my;
use chrono::{DateTime, Utc, TimeZone};
use serde::{Deserialize, Serialize};

use telebot::functions::*;

#[derive(Debug, Serialize, Deserialize)]
struct Quote {
    recnum: u32,
    quote: String,
    author: String,
    date: DateTime<Utc>,
}

impl Quote {
    fn author_nick(&self) -> &str {
        match &self.author.find("!") {
            Some(position) => &self.author[0..*position],
            None => &self.author
        }
    }
}
impl ::std::fmt::Display for Quote {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        write!(f, "{}", format!("{}\n\n-- Quote {} by {} on {}", &self.quote, &self.recnum, &self.author_nick(), &self.date))
    }
}

fn main() {
    let token = env::var("TELEGRAM_BOT_TOKEN").expect("TELEGRAM_BOT_TOKEN not set");
    let db_url = env::var("DIOSTEAMA_DB_URL").expect("DIOSTEAMA_DB_URL not set");
    let pool = my::Pool::new(db_url).expect("Unable to connecto to DB");
    let mut bot = Bot::new(&token).update_interval(200);
    println!("Starting DiosTeama: {}", rquote(&pool));

    let quote = bot.new_cmd("/quote")
    .and_then(move |(bot, msg)| {
        let input = msg.text.unwrap();
        let text = if input.is_empty() {
            rquote(&pool)
        } else {
            match input.parse::<u32>() {
                Ok(qid) => format!("{}", quote_num(&pool, qid)),
                Err(_) =>  format!("{}", quote_like(&pool, &input)),
            }
        };
        // construct a message and return a new future which will be resolved by tokio
        bot.message(msg.chat.id, text).send()
    })
    .for_each(|_| Ok(()));
    bot.run_with(quote);
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
