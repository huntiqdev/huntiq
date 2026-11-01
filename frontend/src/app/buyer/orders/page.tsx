'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { Spinner } from '@/components/ui/Spinner';
import { ShoppingCart, RefreshCw } from 'lucide-react';
import { ORDER_STATUS } from '@/lib/constants/statuses';

const STATUS_NAMES: Record<string, Record<string, string>> = {
  en: { all: 'All', pending: 'Pending', accepted: 'Accepted', completed: 'Completed', declined: 'Declined', cancelled: 'Cancelled' },
  am: { all: 'ሁሉም', pending: 'በመጠባበቅ ላይ', accepted: 'ተቀባይነት አግኝቷል', completed: 'ተጠናቋል', declined: 'ተቀባይነት አላገኘም', cancelled: 'ተሰርዟል' },
  or: { all: 'Hunda', pending: 'Eeggamaa', accepted: 'Fudhatameera', completed: 'Xumurameera', declined: 'Kufaa Ta\'e', cancelled: 'Haqameera' },
  ti: { all: 'ኩሉ', pending: 'ይጽበ ኣሎ', accepted: 'ተቐባልነት ረኺቡ', completed: 'ተወዲኡ', declined: 'ተነጺጉ', cancelled: 'ተሰሪዙ' },
};

export default function BuyerOrdersPage() {
  const { user } = useAuth();
  const { locale } = useLocale();
  const lang = STATUS_NAMES[locale] ? locale : 'en';
  const statusMap = STATUS_NAMES[lang] || STATUS_NAMES.en;

  const [orders, setOrders]   = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter]   = useState('all');

  const load = () => {
    if (!user) return;
    setLoading(true);
    api.getBuyerOrders(user.token)
      .then(o => setOrders(o ?? []))
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, [user]); // eslint-disable-line

  const filtered = filter === 'all' ? orders : orders.filter(o => o.status === filter);
  const statuses = ['all', ...Array.from(new Set(orders.map(o => o.status)))];

  if (loading) return <div className="flex items-center justify-center h-64"><Spinner /></div>;

  const titleText =
    lang === 'am' ? 'የእኔ ጥያቄዎች' : lang === 'or' ? 'Gaaffilee Koo' : lang === 'ti' ? 'ሕቶታተይ' : 'My Requests';
  const subtitleText =
    lang === 'am'
      ? `ከ${orders.length} ውስጥ ${filtered.length} ጥያቄዎች`
      : lang === 'or'
      ? `gaaffilee ${orders.length} keessaa ${filtered.length}`
      : lang === 'ti'
      ? `ካብ ${orders.length} ${filtered.length} ሕቶታት`
      : `${filtered.length} of ${orders.length} requests`;
  const emptyText =
    lang === 'am' ? 'ምንም ጥያቄዎች እስካሁን የሉም' : lang === 'or' ? 'Hamma yoonaatti gaaffiin hin jiru' : lang === 'ti' ? 'ክሳብ ሕጂ ዝቐረበ ሕቶ የለን' : `No ${filter === 'all' ? '' : filter} requests yet`;
  const productLabel =
    lang === 'am' ? 'ምርት' : lang === 'or' ? 'Oomisha' : lang === 'ti' ? 'ፍርያት' : 'Product';
  const offerLabel =
    lang === 'am' ? 'የቀረበ ዋጋ' : lang === 'or' ? 'Gatii Dhiyaate' : lang === 'ti' ? 'ዝቐረበ ዋጋ' : 'Offer';
  const sellerLabel =
    lang === 'am' ? 'ሻጭ' : lang === 'or' ? 'Gurguraa' : lang === 'ti' ? 'ሸያጢ' : 'Seller';

  return (
    <div className="space-y-5 max-w-2xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-black text-zinc-900 dark:text-white">{titleText}</h1>
          <p className="text-xs text-zinc-500 mt-0.5">{subtitleText}</p>
        </div>
        <button onClick={load} className="p-2.5 rounded-xl bg-zinc-100 dark:bg-zinc-800 text-zinc-500 hover:text-zinc-700 dark:hover:text-white transition-all tap-bounce active:scale-90 min-h-[38px] min-w-[38px] flex items-center justify-center" aria-label="Refresh orders">
          <RefreshCw size={14} />
        </button>
      </div>

      {/* Status filter */}
      <div className="flex gap-1.5 flex-wrap">
        {statuses.map(s => (
          <button key={s} onClick={() => setFilter(s)}
            className={`px-3.5 py-2 min-h-[36px] rounded-xl text-xs font-bold capitalize transition-all tap-bounce active:scale-95 ${
              filter === s
                ? 'bg-orange-600 text-white shadow-sm'
                : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500 hover:text-zinc-700 dark:hover:text-white'
            }`}>
            {statusMap[s] || s}
          </button>
        ))}
      </div>

      {filtered.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-48 gap-3 text-zinc-400">
          <ShoppingCart size={32} className="opacity-30" />
          <p className="text-sm">{emptyText}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map((o: any) => (
            <div key={o.id} className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <p className="font-bold text-zinc-900 dark:text-white text-sm">{productLabel} #{o.product_id}</p>
                    <span className={`text-[10px] font-black px-2 py-0.5 rounded-full uppercase ${(ORDER_STATUS[o.status] ?? ORDER_STATUS.pending).color}`}>
                      {statusMap[o.status] || o.status}
                    </span>
                  </div>
                  {o.offer_price && !isNaN(Number(o.offer_price)) && (
                    <p className="text-xs text-zinc-500 mt-1">{offerLabel}: ETB {Number(o.offer_price).toLocaleString()}</p>
                  )}
                  {o.message && (
                    <p className="text-xs text-zinc-400 mt-1 italic">"{o.message}"</p>
                  )}
                  {o.seller_note && (
                    <p className="text-xs text-zinc-500 mt-1 bg-zinc-50 dark:bg-zinc-800 rounded-lg px-2.5 py-1.5">
                      {sellerLabel}: "{o.seller_note}"
                    </p>
                  )}
                </div>
                <p className="text-[10px] text-zinc-400 shrink-0 whitespace-nowrap">
                  {new Date(o.created_at).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: '2-digit' })}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
