import { MongoMemoryServer } from "mongodb-memory-server";
import { connectDB, getDB, getMongoClient } from "../storage/mongo";

let mongoServer: MongoMemoryServer;

beforeAll(async () => {
  process.env.JWT_KEY = "asdfasdf";
  process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";

  mongoServer = await MongoMemoryServer.create();
  const mongoUri = mongoServer.getUri();

  await connectDB(mongoUri);
});

beforeEach(async () => {
  const db = getDB();
  if (db) {
    const collections = await db.collections();
    for (let collection of collections) {
      await collection.deleteMany({});
    }
  }
});

afterAll(async () => {
  const client = getMongoClient();
  if (client) {
    await client.close();
  }
  await mongoServer.stop();
});
