import axios, { AxiosResponse } from "axios";
import request from "supertest";
import { app } from "../../app";

jest.mock("axios");
const mockedAxios = axios as jest.Mocked<typeof axios>;

const signIn = async () => {
  const email = "test@test.com";
  const password = "password";

  const response = await request(app)
    .post("/api/users/signup")
    .send({
      email,
      password,
    })
    .expect(201);

  const cookie = response.get("Set-Cookie");

  if (!cookie) {
    throw new Error("Failed to get cookie from response");
  }
  return cookie;
};

describe("Requests to listings service tests", () => {
  it("returns 401 when no cookie is provided", async () => {
    const response = await request(app).get("/api/listings").send().expect(401);

    expect(response.body.errors[0].message).toMatch(/Not Authorized/i);
  });

  it("forwards a get request and returns successful response", async () => {
    // sign in
    const cookie = await signIn();

    // mock successful response
    const mockedResponse = {
      data: {
        data: [
          {
            id: "id-1",
            type: "listing",
          },
          {
            id: "id-2",
            type: "listing",
          },
        ],
      },
      status: 200,
    };
    mockedAxios.request.mockResolvedValue(mockedResponse);

    // make request with cookie
    const response = await request(app)
      .get("/api/listings")
      .set("Cookie", cookie)
      .send()
      .expect(200);

    expect(response.body).toEqual(mockedResponse.data);
  });

  it("forwards a post request and returns successful response", async () => {
    // sign in
    const cookie = await signIn();

    // mock successful response
    const mockedResponse = {
      data: {
        data: {
          id: "test-id",
          type: "listing",
          title: "test-title",
          description: "description",
        },
      },
      status: 201,
    };
    mockedAxios.request.mockResolvedValue(mockedResponse);

    // make request with cookie
    const response = await request(app)
      .post("/api/listings")
      .set("Cookie", cookie)
      .send({
        data: {
          title: "test-title",
          description: "description",
        },
      })
      .expect(201);

    expect(response.body).toEqual(mockedResponse.data);
  });

  it("forwards a post request and returns error response", async () => {
    // sign in
    const cookie = await signIn();

    // mock successful response
    const mockedResponse = {
      data: {
        data: {
          errors: [
            {
              message: "bad request",
            },
          ],
        },
      },
      status: 400,
    };
    mockedAxios.request.mockResolvedValue(mockedResponse);

    // make request with cookie
    const response = await request(app)
      .post("/api/listings")
      .set("Cookie", cookie)
      .send({
        data: {
          title: "test-title",
          description: "description",
        },
      })
      .expect(400);

    expect(response.body).toEqual(mockedResponse.data);
  });

  it("returns 500 if axios throws an error", async () => {
    const cookie = await signIn();

    mockedAxios.request.mockRejectedValueOnce(new Error("Something broke"));

    const res = await request(app)
      .get("/api/listings/listing-id-1")
      .set("Cookie", cookie)
      .send();

    expect(res.status).toBe(500);
    expect(res.body.errors[0].message).toBe("Internal Error...");
  });
});
