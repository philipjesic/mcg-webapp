import { app } from "./app";
import { connectDB } from "./storage/mongo";

const start = async () => {
  try {
    if (!process.env.JWT_SECRET) {
      throw new Error("JWT_SECRET must be defined");
    }

    if (!process.env.MONGO_URI) {
      throw new Error("MONGO_URI must be defined");
    }

    await connectDB(process.env.MONGO_URI!);
    console.log("Connected to MongoDB...");
  } catch (err) {
    console.error(`Could not run auth service: ${err}`);
    process.exitCode = 1;
    process.exit();
  }
  app.listen(3000, () => {
    console.log("Auth Service listening on port 3000...");
  });
};

start();
