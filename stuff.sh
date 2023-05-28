rm -rf /tmp/lemondrop
>/tmp/lemondrop.exclude
# echo lemondrop_test.go >>/tmp/lemondrop.exclude
mkdir -p /tmp/lemondrop
rg --files --type go | grep -v --file /tmp/lemondrop.exclude | xargs -I{} cp {} /tmp/lemondrop
rg --files | rg go.mod | grep -v --file /tmp/lemondrop.exclude | xargs -I{} cp {} /tmp/lemondrop
rg --files /tmp/lemondrop
txtar-c /tmp/lemondrop >/tmp/lemondrop.txtar
cat /tmp/lemondrop.txtar | pbcopy
