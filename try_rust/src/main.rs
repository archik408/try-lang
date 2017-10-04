extern crate iron;
extern crate postgres;
extern crate router;
extern crate rustc_serialize;
extern crate url;

use iron::*;
use postgres::{Connection, SslMode};

use std::sync::{Arc, Mutex};

mod person_handler;
mod person_dao;

fn init_db(db: &Connection) {
    db.execute(
        concat!(r#"CREATE TABLE IF NOT EXISTS person"#,
                r#"("id" SERIAL PRIMARY KEY, "data" varchar(50),"#,
                r#" "name" varchar(100))"#), &[])
        .unwrap();
}

macro_rules! clone_pass_bound {
    ($arc:ident, $stmt:stmt) => {
        {
            let $arc = $arc.clone();
            $stmt;
        }
    }
}

macro_rules! define_handler {
    ($connection:ident, $router:ident.$method:ident, $route:expr, $handler:path) => {
        clone_pass_bound!(
            $connection,
            $router.$method(
                $route,
                move |req: &mut Request|
                $handler(&*$connection, req),
                ""));
    }
}

macro_rules! define_handlers_family {
    ($connection:ident, $router:ident,
     $( [$method:ident, $route:expr, $handler:path]),+ ) => {
        $( define_handler!($connection, $router.$method, $route, $handler); )+
    }
}

fn serve(db: Connection) {
    let sdb = Arc::new(Mutex::new(db));
    let mut router = router::Router::new();

    define_handlers_family!(
        sdb, router,
        [get, "/api/v1/persons", person_handler::list]
//        [get, "/api/v1/persons/:id", person_handler::get],
//        [post, "/api/v1/persons", person_handler::add],
//        [put, "/api/v1/persons/:id", person_handler::update],
//        [delete, "/api/v1/persons/:id", person_handler::delete]
    );

    Iron::new(router).http("localhost:3000").unwrap();
}

fn main() {
    let db = Connection::connect("postgres://user:pswd@localhost:5432/db", &SslMode::None).unwrap();
    init_db(&db);
    serve(db);
}