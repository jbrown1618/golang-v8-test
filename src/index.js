import { vec, mat, prettyPrint } from "@josh-brown/vector";

const A = mat([
  [0, 1],
  [1, 0],
]);

const x = vec([3, 5]);

global.out = prettyPrint(A.apply(x));
