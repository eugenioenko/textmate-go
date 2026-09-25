def render(name: str, values: list[int]) -> str:
    """Render a compact report.

    The docstring intentionally carries tokenizer state across lines.
    """
    total = sum(values)
    return f"{name}: {total:,}"


print(render("sample", [1, 2, 3]))
