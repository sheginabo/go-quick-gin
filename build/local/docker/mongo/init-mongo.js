// 打印開始日誌
print("Starting MongoDB initialization...");

db = db.getSiblingDB('chatdb');
print("Switched to chatdb...");

// 創建集合
db.createCollection('messages');
print("Created collection: messages");

// 插入初始文檔
db.messages.insertOne({ "init": "initial document" });
print("Inserted initial document into messages collection...");

print("MongoDB initialization script completed.");