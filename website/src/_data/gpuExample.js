const { readFileSync } = require("node:fs");
const path = require("node:path");
module.exports = () => {
  const source = readFileSync(path.resolve(__dirname, "../../..", "examples/gpt2/python/model.py"), "utf8");
  const product = source.match(/^def product\([^\n]+\n(?:[ \t]+[^\n]+\n)+/m)?.[0];
  if (!product) throw new Error("GPT-2 product excerpt moved");
  return product;
};
