interface Box<T extends { id: string }> {
  value: T;
  describe(): string;
}

export const box: Box<{ id: string }> = {
  value: { id: "alpha" },
  describe() {
    return `box:${this.value.id}`;
  },
};
