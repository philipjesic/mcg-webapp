import request from "supertest";
import { app } from "../../app";

describe("User login tests", () => {
  it("responds with a cookie with the correct credentials", async () => {
    await request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(201);

    const response = await request(app)
      .post("/api/users/login")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(201);

    expect(response.get("Set-Cookie")).toBeDefined();
  });

  it("fails when email that does not exist is supplied", async () => {
    await request(app)
      .post("/api/users/login")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(400);
  });

  it("fails when incorrect password is supplied", async () => {
    await request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(201);

    await request(app)
      .post("/api/users/login")
      .send({
        email: "test@test.com",
        password: "wrong password",
      })
      .expect(400);
  });
});
