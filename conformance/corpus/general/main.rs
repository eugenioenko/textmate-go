fn main() {
    /* Nested comments are legal in Rust.
       /* inner */ outer */
    let name = "ferris";
    let message = r#"a raw "quoted" value"#;
    println!("{name}: {message}");
}
