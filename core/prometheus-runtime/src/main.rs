use anyhow::Result;

#[tokio::main]
async fn main() -> Result<()> {
    println!("PROMETHEUS runtime bootstrapped");
    Ok(())
}
