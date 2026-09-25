export const View = ({ title, count }) => (
  <article
    aria-label={
      `${title}: ${count}`
    }
    data-count={count}
  >{title}</article>
);
