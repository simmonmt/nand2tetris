load Max.hack,
output-file Max.out,
compare-to Max.cmp,
output-list RAM[2]%D2.6.2;

set RAM[0] 1,
set RAM[1] 2,
repeat 50 {
  ticktock;
}
output;

set RAM[0] 2,
set RAM[1] 1,
repeat 50 {
  ticktock;
}
output;

set RAM[0] 2,
set RAM[1] 2,
repeat 50 {
  ticktock;
}
output;
