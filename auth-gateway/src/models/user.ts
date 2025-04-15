import { ObjectId } from "mongodb";
import { getDB } from "../storage/mongo";

export interface User {
  _id?: ObjectId;
  email: string;
  passwordHash: string;
  role: "user" | "admin";
}

export const getUserByEmail = async (email: string): Promise<User | null> => {
  const db = getDB();
  return db.collection<User>("users").findOne({ email });
};

export const createUser = async (user: User): Promise<void> => {
  const db = getDB();
  await db.collection<User>("users").insertOne(user);
};
