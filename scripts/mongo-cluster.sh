#!/bin/bash

echo "⏳ 等待 mongo-1 就绪……"
until mongosh --eval 'db.adminCommand({ping:1})' mongodb://mongo-1:27017; do
  sleep 2
done

echo "🚀 执行副本集初始化"
mongosh mongodb://mongo-1:27017 <<'EOF'
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo-1:27017" },
    { _id: 1, host: "mongo-2:27017" },
    { _id: 2, host: "mongo-3:27017" }
  ]
});
EOF

echo "✅ 初始化脚本执行结束"