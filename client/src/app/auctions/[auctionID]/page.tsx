import AuctionDetail from "../../components/AuctionDetail";

interface Props {
  params: {
    listingId: string;
  };
}

export default async function AuctionPage({ params }: Props) {
  /*const res = await fetch(`${process.env.API_BASE_URL}/api/listings/${params.listingId}`, {
    cache: "no-store", // always get latest
  });

  const listing = await res.json();
    */

  const listing = {
    id: "rM2W8MAR",
    title: "2016 Volvo V60 Polestar",
    subtitle: "Only 42,000 miles · AWD · Clean title",
    description: "This V60 Polestar is a rare performance wagon...",
    imageUrl: "https://your-cdn.com/images/v60-front.jpg",
    specs: {
      Mileage: "42,000 miles",
      Transmission: "Automatic",
      Drivetrain: "AWD",
      Color: "Rebel Blue",
      Engine: "3.0L Turbo I6",
    },
    endTime: "2025-05-30T17:00:00.000Z",
    seller: {
      name: "John Doe",
      location: "Vancouver, BC",
    },
  };
  return <AuctionDetail listing={listing} />;
}
