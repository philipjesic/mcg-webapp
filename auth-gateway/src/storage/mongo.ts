import { MongoClient, Db } from "mongodb";

let db: Db;
let client: MongoClient;

export const connectDB = async (uri: string) => {
    client = new MongoClient(uri);
    await client.connect();
    db = client.db();
}

export const getDB = () => {
    if (db) return db;
    else throw new Error("Could not connect to Mongo...");
}

export const getMongoClient = () => {
    if (client) return client;
    else throw new Error('Could not initialize Mongo Client');
}