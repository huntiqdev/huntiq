import { Card } from '@/components/ui/Card';
import { Package, Users, CreditCard, Clock } from 'lucide-react';

interface Props {
  totalProducts: number;
  pendingProducts: number;
  totalSellers: number;
  totalTransactions: number;
}

export function Stats({ totalProducts, pendingProducts, totalSellers, totalTransactions }: Props) {
  const items = [
    { label: 'Total Products', value: totalProducts, icon: Package, color: 'text-blue-600 bg-blue-600/10' },
    { label: 'Pending Review', value: pendingProducts, icon: Clock, color: 'text-yellow-600 bg-yellow-600/10' },
    { label: 'Total Sellers', value: totalSellers, icon: Users, color: 'text-green-600 bg-green-600/10' },
    { label: 'Transactions', value: totalTransactions, icon: CreditCard, color: 'text-orange-600 bg-orange-600/10' },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {items.map(({ label, value, icon: Icon, color }) => (
        <Card key={label}>
          <div className={`w-10 h-10 rounded-xl flex items-center justify-center mb-3 ${color}`}>
            <Icon size={20} />
          </div>
          <p className="text-2xl font-black">{value}</p>
          <p className="text-xs font-bold opacity-50 mt-1">{label}</p>
        </Card>
      ))}
    </div>
  );
}
