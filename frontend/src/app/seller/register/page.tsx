import { redirect } from 'next/navigation';

export default function SellerRegisterPage() {
  redirect('/auth?as=seller&tab=register');
}
