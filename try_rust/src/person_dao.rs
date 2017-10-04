use postgres::Connection;

use std::sync::Mutex;

pub fn insert(db: &Connection, name: &str, data: &str) -> ::postgres::Result<u64> {
    db.execute("INSERT INTO person VALUES (default, $1, $2)", &[&name, &data])
}

pub fn remove(db: &Connection, ids: &[i32]) -> ::postgres::Result<u64> {
    let stmt = db.prepare("DELETE FROM person WHERE id=$1")
        .unwrap();
    for id in ids {
        try!(stmt.execute(&[id]));
    }
    Ok(0)
}

pub fn update(db: &Connection, id: i32, name: &str, data: &str) -> ::postgres::Result<()> {
    let tx: ::postgres::Transaction = db.transaction()
        .unwrap();
    tx.execute("UPDATE person SET name = $1, data = $2 WHERE id = $3", &[&name, &data, &id])
        .unwrap();
    tx.set_commit();
    tx.finish()
}

pub fn show(db: &Connection, arg: Option<&str>) -> ::postgres::Result<Vec<Person>> {
    let s = match arg {
        Some(s) => format!("WHERE name LIKE '%{}%'", s),
        None => "".to_owned(),
    };
    let stmt = db.prepare(&format!("SELECT * FROM person {} ORDER BY id", s))
        .unwrap();
    let rows = stmt.query(&[]).unwrap();
    let size = rows.iter().count();
    let mut results = Vec::with_capacity(size);
    for row in rows {
        let record = Person {
            id: row.get("id"),
            name: row.get("name"),
            data: row.get("data"),
        };
        results.push(record)
    }
    Ok(results)
}

#[derive(RustcEncodable, RustcDecodable)]
pub struct Person {
    id: Option<i32>,
    pub name: String,
    pub data: String,
}

pub fn read(sdb: &Mutex<Connection>, name: Option<&str>) -> Result<Vec<Person>, ()> {
    if let Ok(rs) = show(&*sdb.lock().unwrap(), name) {
        Ok(rs)
    } else {
        Err(())
    }
}

pub fn read_one(sdb: &Mutex<Connection>, id: i32) -> Result<Person, ()> {
    let db = &*sdb.lock().unwrap();
    let stmt = db.prepare("SELECT * FROM person WHERE id = $1")
        .unwrap();
    if let Ok(rows) = stmt.query(&[&id]) {
        let mut iter = rows.iter();
        if iter.len() != 1 {
            return Err(());
        }
        let row = iter.next().unwrap();
        let record = Person {
            id: row.get("id"),
            name: row.get("name"),
            data: row.get("data"),
        };

        Ok(record)
    } else {
        Err(())
    }
}