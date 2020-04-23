{ s = substr($0, 1, 1);
  val = 0
  if (s == "+") { val = 1 }
  else if (s == "-") { val = -1 }
  if (val == 0 || substr($0, 2, 1) == s) {
    next
  }
  printf("Got a %s in %s\n", s, $0);
  s = substr($0, 2);
  counts[s] += val
}
END  { status = 0;
       for (idx in counts) {
         if (counts[idx] != 0) {
            printf("Unmatched: %s: %d\n", idx, counts[idx]);
            status = 1
         }
       }
       exit status
}
