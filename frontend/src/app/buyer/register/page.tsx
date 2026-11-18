import { redirect } from 'next/navigation';

export default function BuyerRegisterPage() {
  redirect('/auth?as=buyer&tab=register');
}
