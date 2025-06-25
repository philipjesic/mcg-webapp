export interface Bid {
  auction_id: string;
  user_id: string;
  amount: number;
  timestamp: string;
}

export interface BidRequest {
  data: Bid;
}
