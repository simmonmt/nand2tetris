load Add.hack,
output-file Add.out,
compare-to Add.cmp,
output-list RAM[0]%D2.6.2;

repeat 50 {
  ticktock;
}
output;
