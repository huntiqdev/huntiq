import { Card } from '@/components/ui/Card';
import { Package, TrendingUp, Clock, CheckCircle } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';

interface Stats {
  total: number;
  pending: number;
  approved: number;
  totalRevenue: number;
}

interface Props { stats: Stats }

export function SellerDashboard({ stats }: Props) {
  const { t } = useLocale();

  const items = [
    { label: t.ui.my_products,        value: stats.total,                              icon: Package,      color: 'text-blue-600 bg-blue-600/10' },
    { label: t.ui.pending_products,   value: stats.pending,                            icon: Clock,        color: 'text-yellow-600 bg-yellow-600/10' },
    { label: t.ui.approve,           value: stats.approved,                           icon: CheckCircle,  color: 'text-green-600 bg-green-600/10' },
    { label: t.ui.total_earnings,     value: `ETB ${stats.totalRevenue.toLocaleString()}`, icon: TrendingUp, color: 'text-orange-600 bg-orange-600/10' },
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
