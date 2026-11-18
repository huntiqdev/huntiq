import { redirect } from 'next/navigation';

export default function RentPage() {
  redirect('/products?type=rent');
}
