import AuctionDetail from "../../components/AuctionDetail";

interface Props {
  params: {
    auctionID: string;
  };
}

export default async function AuctionPage({ params }: Props) {
  const { auctionID } = await params;
  const listingsAddr = process.env.LISTINGS_SERVICE || "";
  const listingsPort = process.env.LISTING_SERVICES_PORT || 3000;
  console.log(
    `http://${listingsAddr}:${listingsPort}/api/listings/${auctionID}`
  );
  const res = await fetch(
    `http://${listingsAddr}:${listingsPort}/api/listings/${auctionID}`
  );

  console.log(res);
  const listingsResponse = await res.json();
  const listing = listingsResponse.data[0];

  return <AuctionDetail listing={listing} />;
}
