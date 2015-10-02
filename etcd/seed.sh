#!/bin/bash

. ../shared.sh

# "%s:%s@tcp(%s:%s)/%s?parseTime=true&strict=true&sql_notes=false"
read -d '' JSON <<EOF || true
[ { "name": "alpha"
  , "driver": "mysql"
  , "dsn": "root:chaostestingrootpassword@tcp(172.17.42.1:3306)/chaostesting?parseTime=true&strict=true&sql_notes=false&timeout=3s&read_timeout=3s&write_timeout=3s"
  , "weight": 100 }
, { "name": "beta"
  , "driver": "mysql"
  , "dsn": "root:chaostestingrootpassword@tcp(172.17.42.1:3307)/chaostesting?parseTime=true&strict=true&sql_notes=false&timeout=3s&read_timeout=3s&write_timeout=3s"
  , "weight": 100 }
]
EOF

./ctl.sh set /chaostesting/datasources "$JSON"
