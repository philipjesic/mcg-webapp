import request from "supertest";
import { app } from "../../app";

describe("User Signup Tests", () => {
  it("should return 201 and sets a cookie after sucessful sign up", async () => {
    const response = await request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(201);
    expect(response.get("Set-Cookie")).toBeDefined();
  });

  it("should return 400 for missing username and password", async () => {
    await request(app)
      .post("/api/users/signup")
      .send({
        password: "password123",
      })
      .expect(400);

    await request(app)
      .post("/api/users/signup")
      .send({
        email: "email@mail.com",
      })
      .expect(400);
  });

  it("returns a 400 with an invalid email", async () => {
    return request(app)
      .post("/api/users/signup")
      .send({
        email: "alskdflaskjfd",
        password: "password",
      })
      .expect(400);
  });

  it("returns a 400 with an invalid password", async () => {
    return request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "p",
      })
      .expect(400);
  });

  it("disallows duplicate emails", async () => {
    const firstRes = await request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(201);

    const response = await request(app)
      .post("/api/users/signup")
      .send({
        email: "test@test.com",
        password: "password",
      })
      .expect(400);
  });
});
