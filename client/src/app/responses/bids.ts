export interface Bid {
  id: string;
  type: string;
  auction_id: string;
  user_id: string;
  amount: number;
  timestamp: string;
}

export interface BidResponse {
  data: Bid;
}
