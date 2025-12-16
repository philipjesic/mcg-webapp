export interface Listing {
  id: string;
  type: string;
  title: string;
  description: string;
  location: string;
  endtime: string;
  seller: {
    name: string;
    location: string;
    id: string;
  },
  specifications: {
    make: string;
    model: string;
    engine: string;
    colour: string;
    transmission: string;
  }
}

export interface ListingsResponse {
  data: Listing[];
}
