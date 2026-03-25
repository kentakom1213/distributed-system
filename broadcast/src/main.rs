use clap::Parser;

/// Initialize broadcast node
#[derive(Debug, clap::Parser)]
struct Args {
    #[arg(long, required = true)]
    id: u16,
    #[arg(short, long, required = true)]
    port: u16,
}

fn main() {
    let args = Args::parse();

    println!("{:?}", args);
}
