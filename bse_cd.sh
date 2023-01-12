marketview|tr ',' '\t'|fzf -d '\t' --with-nth 2,3,4,5 --exact| tr '\t' ',' |cut -d',' -f2
