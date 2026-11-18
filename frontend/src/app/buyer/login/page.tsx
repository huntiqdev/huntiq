import { redirect } from 'next/navigation';

export default function BuyerLoginPage() {
  redirect('/auth?as=buyer');
}
