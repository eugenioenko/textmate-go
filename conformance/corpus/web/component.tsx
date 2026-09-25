type CardProps = {
  title: string;
  count: number;
};

export function Card({ title, count }: CardProps) {
  /* The expression and quoted attribute both span physical lines. */
  const label = `items: ${count}`;
  return (
    <section
      className="card"
      aria-label={
        count > 0
          ? label
          : "empty"
      }
      data-title={title}
    >
      <strong>{title}</strong>
    </section>
  );
}
