export interface Listing {
  id: string;
  title: string;
  description: string;
}

export interface ListingResponse {
  data: Listing[];
}
