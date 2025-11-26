/// This is a semi advanced example of a simple economy plugin.
/// we are gonna use sqlite, to store user money.
/// two commands:
/// pay: pay yourself money
/// bal: view your balance / money
use dragonfly_plugin::{
    Plugin, PluginRunner, Server,
    event::{EventContext, EventHandler},
    event_handler, types,
};
use sqlx::{SqlitePool, sqlite::SqlitePoolOptions};

#[derive(Plugin)]
#[plugin(
    id = "rustic-economy",
    name = "Rustic Economy",
    version = "0.1.0",
    api = "1.0.0"
)]
struct RusticEconomy {
    db: SqlitePool,
}

/// This impl is just a helper for dealing with our SQL stuff.
impl RusticEconomy {
    async fn new() -> Result<Self, Box<dyn std::error::Error>> {
        // Create database connection
        let db = SqlitePoolOptions::new()
            .max_connections(5)
            .connect("sqlite:economy.db")
            .await?;

        // Create table if it doesn't exist
        sqlx::query(
            "CREATE TABLE IF NOT EXISTS users (
                uuid TEXT PRIMARY KEY,
                balance REAL NOT NULL DEFAULT 0.0
            )",
        )
        .execute(&db)
        .await?;

        Ok(Self { db })
    }

    async fn get_balance(&self, uuid: &str) -> Result<f64, sqlx::Error> {
        let result: Option<(f64,)> = sqlx::query_as("SELECT balance FROM users WHERE uuid = ?")
            .bind(uuid)
            .fetch_optional(&self.db)
            .await?;

        Ok(result.map(|(bal,)| bal).unwrap_or(0.0))
    }

    async fn add_money(&self, uuid: &str, amount: f64) -> Result<f64, sqlx::Error> {
        // Insert or update user balance
        sqlx::query(
            "INSERT INTO users (uuid, balance) VALUES (?, ?)
             ON CONFLICT(uuid) DO UPDATE SET balance = balance + ?",
        )
        .bind(uuid)
        .bind(amount)
        .bind(amount)
        .execute(&self.db)
        .await?;

        self.get_balance(uuid).await
    }
}

#[event_handler]
impl EventHandler for RusticEconomy {
    async fn on_chat(&self, server: &Server, event: &mut EventContext<'_, types::ChatEvent>) {
        let message = &event.data.message;
        let player_uuid = &event.data.player_uuid;

        // Handle commands
        if message.starts_with("!pay") {
            // ALWAYS RESPOND TO THE SERVER ASAP.
            // THEN DO SIDE WORK THAT MIGHT TAKE A GOOD BIT.
            let parts: Vec<&str> = message.split_whitespace().collect();
            if parts.len() != 2 {
                server
                    .send_chat(player_uuid.clone(), "Usage: !pay <amount>".to_string())
                    .await
                    .expect("Bad error handling womp.");
                return;
            }

            let amount: f64 = match parts[1].parse() {
                Ok(amt) if amt > 0.0 => amt,
                _ => {
                    server
                        .send_chat(
                            player_uuid.clone(),
                            "Please provide a valid positive amount!".to_string(),
                        )
                        .await
                        .expect("Bad error handling sad.");
                    return;
                }
            };

            match self.add_money(player_uuid, amount).await {
                Ok(new_balance) => {
                    server
                        .send_chat(
                            player_uuid.clone(),
                            format!("Added ${:.2}! New balance: ${:.2}", amount, new_balance),
                        )
                        .await
                        .expect("again error handling is bad");
                }
                Err(e) => {
                    eprintln!("Database error: {}", e);
                    server
                        .send_chat(player_uuid.clone(), "Error processing payment!".to_string())
                        .await
                        .expect("Bad error handling");
                }
            }
            event.cancel().await;
            println!("I HAVE SENT THE CANCEL FULLY");
        } else if message.starts_with("!bal") {
            // YET AGAIN RESPOND ASAP TO SERVER.
            event.cancel().await;
            match self.get_balance(player_uuid).await {
                Ok(balance) => {
                    // These are fine to happen later.
                    // as they aren't bounded by the timout of server responses.
                    server
                        .send_chat(
                            player_uuid.clone(),
                            format!("Your balance: ${:.2}", balance),
                        )
                        .await
                        .expect("oh shit i need to handle this better.");
                }
                Err(e) => {
                    eprintln!("Database error: {}", e);
                    server
                        .send_chat(player_uuid.clone(), "Error checking balance!".to_string())
                        .await
                        .expect("again bad error handling");
                }
            }
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Starting the plugin...");
    println!("Initializing database...");

    let plugin = RusticEconomy::new().await?;

    PluginRunner::run(plugin, "tcp://127.0.0.1:50050").await
}
