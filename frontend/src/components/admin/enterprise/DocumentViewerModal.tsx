'use client';

import React, { useState, useEffect } from 'react';
import { 
  X, FileText, Download, ShieldCheck, ZoomIn, ZoomOut, Printer, 
  Eye, CheckCircle2, Lock, Fingerprint, Copy, Check, MessageSquare, 
  AlertTriangle, FileCode, ShieldAlert, Sparkles, Bookmark, Info
} from 'lucide-react';
import { DataRowItem } from './EnterpriseDataTable';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';

interface DocumentViewerModalProps {
  isOpen: boolean;
  documentName: string | null;
  selectedRow: DataRowItem | null;
  onClose: () => void;
  onStatusUpdated?: () => void;
}

interface CommentItem {
  id: string;
  author: string;
  role: string;
  text: string;
  time: string;
}

export function DocumentViewerModal({ isOpen, documentName, selectedRow, onClose, onStatusUpdated }: DocumentViewerModalProps) {
  const { user } = useAuth();
  const [zoomScale, setZoomScale] = useState<number>(1.0);
  const [activeTab, setActiveTab] = useState<'compliance' | 'comments' | 'crypto'>('compliance');
  const [filterMode, setFilterMode] = useState<'original' | 'high-contrast' | 'sepia'>('original');
  const [showWatermark, setShowWatermark] = useState<boolean>(true);
  const [selectedClause, setSelectedClause] = useState<number | null>(null);
  
  // Checklist State
  const [checklist, setChecklist] = useState({
    identity: false,
    escrow: false,
    permit: false,
    aml: false,
  });

  // Comments State
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [newComment, setNewComment] = useState<string>('');
  const [copiedHash, setCopiedHash] = useState<boolean>(false);

  // Success/Notification Toast
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'info' | 'error' } | null>(null);

  // Set default comments when row changes
  useEffect(() => {
    if (selectedRow) {
      setComments([
        {
          id: 'cmt-1',
          author: 'Alex Mercer',
          role: 'Lead Security Auditor',
          text: `Telemetry logs confirm IP matching registry details for ID ${selectedRow.id}.`,
          time: '2 hours ago'
        },
        {
          id: 'cmt-2',
          author: 'Sophia Chen',
          role: 'Escrow Officer',
          text: `Escrow voucher verification token matches bank clearance code for $${selectedRow.specDetails.escrowAmount?.toLocaleString() || 'N/A'}.`,
          time: '30 mins ago'
        }
      ]);
      // Reset checklist
      setChecklist({
        identity: selectedRow.status === 'Approved',
        escrow: selectedRow.status === 'Approved' || selectedRow.status === 'In Escrow',
        permit: selectedRow.status === 'Approved',
        aml: selectedRow.status === 'Approved'
      });
    }
  }, [selectedRow]);

  if (!isOpen || !documentName || !selectedRow) return null;

  // Calculate Verification Progress
  const checkedCount = Object.values(checklist).filter(Boolean).length;
  const progressPercent = Math.round((checkedCount / 4) * 100);

  const handleToggleCheck = (key: keyof typeof checklist) => {
    setChecklist(prev => ({ ...prev, [key]: !prev[key] }));
  };

  const handleAddComment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    
    const commentObj: CommentItem = {
      id: `cmt-${Date.now()}`,
      author: 'You (Current Admin)',
      role: 'Super Administrator',
      text: newComment.trim(),
      time: 'Just now'
    };

    setComments(prev => [...prev, commentObj]);
    setNewComment('');
    showToast('Annotation added to compliance audit trail.', 'success');
  };

  const handleCopyHash = () => {
    const fakeHash = `sha256-8F9083BA571E9023AC${selectedRow.id}E47F20E90312FDC8B95A2`;
    navigator.clipboard.writeText(fakeHash);
    setCopiedHash(true);
    showToast('Cryptographic hash copied to clipboard.', 'info');
    setTimeout(() => setCopiedHash(false), 2000);
  };

  const showToast = (message: string, type: 'success' | 'info' | 'error') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  const handleAuditAction = (status: 'Approved' | 'Rejected' | 'Suspended') => {
    if (status === 'Approved' && progressPercent < 100) {
      showToast('Cannot approve. Please verify all checklist items first.', 'error');
      return;
    }
    
    if (!selectedRow) return;
    const numericId = parseInt(selectedRow.id.split('-')[1]);
    if (isNaN(numericId)) return;

    const backendStatus = status === 'Approved' ? 'approved' : 'rejected';
    
    if (user?.token) {
      showToast('Submitting moderation status to backend...', 'info');
      api.adminReviewProduct(
        numericId,
        {
          status: backendStatus,
          notes: `Reviewed and certified via Document Verification Workspace (${documentName})`
        },
        user.token
      )
        .then(() => {
          showToast(`Document status successfully marked as: ${status}`, 'success');
          if (onStatusUpdated) {
            onStatusUpdated();
          }
        })
        .catch((err: any) => {
          console.error('Failed to review product:', err);
          showToast(`Validation check: ${err.message || 'Validation error'}`, 'error');
        });
    } else {
      showToast(`Document status successfully marked as: ${status}`, 'success');
    }
  };

  // Document formatting styles based on filter modes
  let documentClass = "bg-zinc-900 border border-zinc-800 text-zinc-100";
  if (filterMode === 'high-contrast') {
    documentClass = "bg-white border-2 border-black text-black font-mono";
  } else if (filterMode === 'sepia') {
    documentClass = "bg-[#f4ecd8] border border-[#e4dcc8] text-[#5b4636] font-serif";
  }

  // Simulated Document hash
  const docHash = `SHA256: 0x8F90...A312-${selectedRow.id}`;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/90 backdrop-blur-md transition-opacity">
      <div className="w-full max-w-6xl bg-zinc-950 border border-zinc-800 rounded-3xl shadow-2xl overflow-hidden flex flex-col h-[90vh] animate-in zoom-in-95 duration-200">
        
        {/* Top Header Bar */}
        <div className="p-4 bg-zinc-900/90 border-b border-zinc-800/80 flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
              <FileText size={20} />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h3 className="text-sm font-extrabold text-white">
                  {documentName}
                </h3>
                <span className="px-2 py-0.5 text-[10px] font-bold font-mono bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 rounded-md">
                  Digitally Sealed PDF
                </span>
              </div>
              <p className="text-xs text-zinc-400 mt-0.5">
                Record Ref: <span className="font-mono text-zinc-300 font-bold">{selectedRow.id}</span> ({selectedRow.title})
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            {/* Viewport Control Actions */}
            <div className="flex items-center bg-zinc-950 border border-zinc-800 rounded-xl p-1 text-zinc-400 gap-1 mr-2">
              <button 
                onClick={() => setZoomScale(z => Math.max(z - 0.1, 0.8))}
                className="p-1.5 hover:bg-zinc-800 hover:text-white rounded-lg transition-colors"
                title="Zoom Out"
              >
                <ZoomOut size={14} />
              </button>
              <span className="text-[10px] px-1 font-mono font-bold">{Math.round(zoomScale * 100)}%</span>
              <button 
                onClick={() => setZoomScale(z => Math.min(z + 0.1, 1.5))}
                className="p-1.5 hover:bg-zinc-800 hover:text-white rounded-lg transition-colors"
                title="Zoom In"
              >
                <ZoomIn size={14} />
              </button>
            </div>

            <button
              onClick={() => {
                window.print();
                showToast('Spooling document to printer...', 'info');
              }}
              className="p-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5"
              title="Print Document"
            >
              <Printer size={14} />
            </button>

            <button
              onClick={() => showToast(`Simulated download of ${documentName}`, 'success')}
              className="px-3.5 py-2 bg-orange-600 hover:bg-orange-500 text-white rounded-xl text-xs font-extrabold transition-all flex items-center gap-1.5 shadow-lg shadow-orange-600/10"
            >
              <Download size={14} /> Download File
            </button>

            <button 
              onClick={onClose} 
              className="p-2 text-zinc-400 hover:text-white rounded-xl hover:bg-zinc-800 border border-transparent hover:border-zinc-700 transition-all ml-2"
            >
              <X size={20} />
            </button>
          </div>
        </div>

        {/* Toast Notification */}
        {toast && (
          <div className={`mx-6 mt-4 p-3 rounded-xl border flex items-center justify-between text-xs font-bold animate-in slide-in-from-top-2 duration-200 z-50 ${
            toast.type === 'success' 
              ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400' 
              : toast.type === 'error'
              ? 'bg-red-500/10 border-red-500/30 text-red-400'
              : 'bg-blue-500/10 border-blue-500/30 text-blue-400'
          }`}>
            <span className="flex items-center gap-2">
              {toast.type === 'success' ? <CheckCircle2 size={15} /> : <AlertTriangle size={15} />}
              {toast.message}
            </span>
            <button onClick={() => setToast(null)} className="hover:text-white">✕</button>
          </div>
        )}

        {/* Main Work Area: Dual Pane Layout */}
        <div className="flex-1 overflow-hidden flex flex-col md:flex-row">
          
          {/* Left Pane: Interactive Document Canvas Viewport */}
          <div className="flex-1 bg-zinc-950 p-4 lg:p-6 overflow-y-auto border-r border-zinc-900 flex flex-col items-center justify-start space-y-4">
            
            {/* Filter & Watermark Controls */}
            <div className="w-full max-w-2xl flex items-center justify-between text-xs text-zinc-400">
              <div className="flex gap-1.5 items-center">
                <span>View Filters:</span>
                {(['original', 'high-contrast', 'sepia'] as const).map((mode) => (
                  <button
                    key={mode}
                    onClick={() => setFilterMode(mode)}
                    className={`px-2.5 py-1 rounded-lg capitalize border font-bold ${
                      filterMode === mode 
                        ? 'bg-zinc-800 border-zinc-700 text-white' 
                        : 'bg-zinc-900/40 border-zinc-800/80 hover:text-white'
                    }`}
                  >
                    {mode.replace('-', ' ')}
                  </button>
                ))}
              </div>

              <label className="flex items-center gap-1.5 cursor-pointer font-bold select-none">
                <input 
                  type="checkbox" 
                  checked={showWatermark}
                  onChange={() => setShowWatermark(!showWatermark)}
                  className="rounded border-zinc-800 text-orange-500 focus:ring-0 focus:ring-offset-0 bg-zinc-900"
                />
                <span>Holographic Seals</span>
              </label>
            </div>

            {/* Document Visual Layout Container */}
            <div 
              className={`w-full max-w-2xl rounded-2xl p-8 shadow-2xl relative transition-all duration-300 ${documentClass}`}
              style={{ 
                transform: `scale(${zoomScale})`, 
                transformOrigin: 'top center',
                transition: 'transform 0.15s ease-out'
              }}
            >
              {/* Official Stamp Overlay */}
              {showWatermark && (
                <div className="absolute top-8 right-8 border-2 border-emerald-500/60 text-emerald-500 px-4 py-2.5 rounded-xl uppercase font-sans font-black tracking-widest text-[10px] rotate-12 bg-emerald-950/20 flex items-center gap-1.5 shadow-lg select-none pointer-events-none transition-opacity duration-200">
                  <ShieldCheck size={14} /> Official Escrow Vault Seal
                </div>
              )}

              {/* Watermark Diagonal Text */}
              {showWatermark && (
                <div className="absolute inset-0 flex items-center justify-center pointer-events-none select-none overflow-hidden opacity-5 z-0">
                  <div className="text-zinc-500 font-extrabold text-5xl font-mono uppercase tracking-[0.2em] -rotate-45">
                    SELLBUY COMPLIANCE
                  </div>
                </div>
              )}

              {/* Document Header Section */}
              <div className={`border-b pb-4 mb-6 space-y-1 z-10 relative ${
                filterMode === 'high-contrast' ? 'border-black' : 'border-zinc-800'
              }`}>
                <div className={`text-[10px] font-mono uppercase tracking-widest font-extrabold ${
                  filterMode === 'high-contrast' ? 'text-black' : 'text-orange-400'
                }`}>
                  SELLBUY ENTERPRISE SAAS COMPLIANCE LEDGER
                </div>
                <h2 className="text-lg font-black tracking-tight font-serif">
                  AGREEMENT OF LEASE & ESCROW LOCK
                </h2>
                <div className="text-[10px] font-mono opacity-80">{docHash} • Verified on {selectedRow.dateCreated}</div>
              </div>

              {/* Interactive Clauses Body */}
              <div className="space-y-4 font-sans text-xs leading-relaxed z-10 relative">
                
                {/* Clause 1 */}
                <div 
                  onClick={() => setSelectedClause(1)}
                  className={`p-4 rounded-xl border transition-all cursor-pointer ${
                    selectedClause === 1 
                      ? 'bg-orange-500/10 border-orange-500/50 ring-1 ring-orange-500/30' 
                      : filterMode === 'original'
                      ? 'bg-zinc-950/50 border-zinc-800/80 hover:border-zinc-700/80 hover:bg-zinc-900/30'
                      : 'hover:opacity-90 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between font-extrabold uppercase tracking-wider text-[10px] opacity-75 mb-1">
                    <span>Section 1: Contracting Parties & Asset Allocation</span>
                    {selectedClause === 1 && <span className="text-orange-400 font-bold">Inspecting</span>}
                  </div>
                  <p>
                    This Lease Agreement is executed between <strong className="font-extrabold">{selectedRow.specDetails.locationOrSeller}</strong> (hereinafter referred to as the Lessor / Authorized Entity) and the verified tenant on record (Lessees). The asset listed is verified under identifier <strong className="font-mono">{selectedRow.id}</strong>.
                  </p>
                </div>

                {/* Clause 2 */}
                <div 
                  onClick={() => setSelectedClause(2)}
                  className={`p-4 rounded-xl border transition-all cursor-pointer ${
                    selectedClause === 2 
                      ? 'bg-orange-500/10 border-orange-500/50 ring-1 ring-orange-500/30' 
                      : filterMode === 'original'
                      ? 'bg-zinc-950/50 border-zinc-800/80 hover:border-zinc-700/80 hover:bg-zinc-900/30'
                      : 'hover:opacity-90 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between font-extrabold uppercase tracking-wider text-[10px] opacity-75 mb-1">
                    <span>Section 2: Rental Valuation & Escrow Escort Lock</span>
                    {selectedClause === 2 && <span className="text-orange-400 font-bold">Inspecting</span>}
                  </div>
                  <p>
                    The agreed rental valuation stands at <strong className="text-orange-400 font-mono">${selectedRow.price.toLocaleString()}</strong>{selectedRow.pricePeriod || '/month'}. 
                    A deposit of <strong className="text-emerald-400 font-mono">${selectedRow.specDetails.escrowAmount?.toLocaleString() || '3,600'}</strong> is verified as securely locked inside the **SellBuy Unified Escrow Vault**. This balance cannot be liquidated or disbursed without dual cryptographic signatures from the Lessor and Super Administrator.
                  </p>
                </div>

                {/* Clause 3 */}
                <div 
                  onClick={() => setSelectedClause(3)}
                  className={`p-4 rounded-xl border transition-all cursor-pointer ${
                    selectedClause === 3 
                      ? 'bg-orange-500/10 border-orange-500/50 ring-1 ring-orange-500/30' 
                      : filterMode === 'original'
                      ? 'bg-zinc-950/50 border-zinc-800/80 hover:border-zinc-700/80 hover:bg-zinc-900/30'
                      : 'hover:opacity-90 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between font-extrabold uppercase tracking-wider text-[10px] opacity-75 mb-1">
                    <span>Section 3: Maintenance Protocols & Operational SLA</span>
                    {selectedClause === 3 && <span className="text-orange-400 font-bold">Inspecting</span>}
                  </div>
                  <p>
                    The maintenance records specify: <span className="font-semibold">{selectedRow.specDetails.maintenanceHistory || 'Regular scheduled checkups apply.'}</span>. 
                    Lessor warrants that all equipment, HVAC systems, plumbing, and structural items remain active and compliant with safety guidelines.
                  </p>
                </div>

                {/* Clause 4 */}
                <div 
                  onClick={() => setSelectedClause(4)}
                  className={`p-4 rounded-xl border transition-all cursor-pointer ${
                    selectedClause === 4 
                      ? 'bg-orange-500/10 border-orange-500/50 ring-1 ring-orange-500/30' 
                      : filterMode === 'original'
                      ? 'bg-zinc-950/50 border-zinc-800/80 hover:border-zinc-700/80 hover:bg-zinc-900/30'
                      : 'hover:opacity-90 border-transparent'
                  }`}
                >
                  <div className="flex items-center justify-between font-extrabold uppercase tracking-wider text-[10px] opacity-75 mb-1">
                    <span>Section 4: Signatures & Digital Audits</span>
                    {selectedClause === 4 && <span className="text-orange-400 font-bold">Inspecting</span>}
                  </div>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
                    <div className={`p-2.5 rounded-lg border text-[11px] ${
                      filterMode === 'original' ? 'bg-zinc-900 border-zinc-800' : 'bg-transparent border-current'
                    }`}>
                      <div className="opacity-60">Lessor Identification Status:</div>
                      <div className="text-emerald-500 font-bold mt-0.5 flex items-center gap-1">
                        <ShieldCheck size={12} /> ID Verified Match
                      </div>
                    </div>
                    <div className={`p-2.5 rounded-lg border text-[11px] ${
                      filterMode === 'original' ? 'bg-zinc-900 border-zinc-800' : 'bg-transparent border-current'
                    }`}>
                      <div className="opacity-60">Crypto Signature Clearance:</div>
                      <div className="text-emerald-500 font-bold mt-0.5 flex items-center gap-1">
                        <Lock size={12} /> Sealed & Locked
                      </div>
                    </div>
                  </div>
                </div>

              </div>

              {/* Document Inspector Footer */}
              <div className={`pt-4 mt-6 border-t flex items-center justify-between text-[9px] font-mono opacity-65 ${
                filterMode === 'high-contrast' ? 'border-black' : 'border-zinc-800'
              }`}>
                <span>SellBuy Compliance Engine v2.8</span>
                <span>Security Level: Tier 3 Audited</span>
              </div>
            </div>
          </div>

          {/* Right Pane: Interactive Compliance & Security Workspace */}
          <div className="w-full md:w-96 bg-zinc-900 border-t md:border-t-0 border-zinc-800/80 flex flex-col shrink-0">
            
            {/* Right Pane Navigation Tabs */}
            <div className="flex border-b border-zinc-800 bg-zinc-950/80 p-2 gap-1.5">
              {(['compliance', 'comments', 'crypto'] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`flex-1 py-2 rounded-xl text-center font-extrabold text-[10px] uppercase tracking-wider transition-all ${
                    activeTab === tab 
                      ? 'bg-orange-600 text-white shadow-md' 
                      : 'text-zinc-400 hover:text-white hover:bg-zinc-900'
                  }`}
                >
                  {tab}
                </button>
              ))}
            </div>

            {/* Scrollable Tab Contents */}
            <div className="flex-1 overflow-y-auto p-5 space-y-6">
              
              {/* Tab 1: Compliance Auditing Checklist */}
              {activeTab === 'compliance' && (
                <div className="space-y-5">
                  <div className="space-y-1">
                    <h4 className="text-xs font-black text-white flex items-center gap-1.5 uppercase tracking-wide">
                      <ShieldCheck size={15} className="text-orange-500" /> Compliance Audit Matrix
                    </h4>
                    <p className="text-[11px] text-zinc-400">Perform itemized visual review to calculate verification trust score.</p>
                  </div>

                  {/* Trust Score Progress Ring/Bar */}
                  <div className="p-4 bg-zinc-950 rounded-2xl border border-zinc-800 space-y-3">
                    <div className="flex items-center justify-between text-xs font-bold">
                      <span className="text-zinc-300">Auditing Verification Progress</span>
                      <span className={`${
                        progressPercent === 100 ? 'text-emerald-400' : 'text-orange-400'
                      } font-mono`}>{progressPercent}%</span>
                    </div>

                    <div className="w-full bg-zinc-900 rounded-full h-2 overflow-hidden border border-zinc-800/50">
                      <div 
                        className={`h-full rounded-full transition-all duration-300 ${
                          progressPercent === 100 
                            ? 'bg-gradient-to-r from-emerald-500 to-teal-400 shadow-lg shadow-emerald-500/20' 
                            : 'bg-gradient-to-r from-orange-500 to-amber-400'
                        }`} 
                        style={{ width: `${progressPercent}%` }}
                      />
                    </div>

                    {progressPercent === 100 ? (
                      <div className="text-[10px] text-emerald-400 flex items-center gap-1.5 font-semibold bg-emerald-500/5 p-1.5 rounded-lg border border-emerald-500/10">
                        <CheckCircle2 size={12} /> All protocols checked. Secure to authorize listing.
                      </div>
                    ) : (
                      <div className="text-[10px] text-zinc-400 flex items-center gap-1.5 font-medium">
                        <Info size={12} className="text-orange-400 shrink-0" /> Pending review completion.
                      </div>
                    )}
                  </div>

                  {/* Checklist Items */}
                  <div className="space-y-2.5">
                    {/* Check 1 */}
                    <button
                      onClick={() => handleToggleCheck('identity')}
                      className={`w-full p-3 rounded-xl border text-left flex items-start gap-3 transition-all ${
                        checklist.identity 
                          ? 'bg-emerald-500/5 border-emerald-500/30 text-emerald-400' 
                          : 'bg-zinc-950/40 border-zinc-800 hover:border-zinc-700/80 text-zinc-300'
                      }`}
                    >
                      <input 
                        type="checkbox" 
                        checked={checklist.identity}
                        readOnly
                        className="rounded border-zinc-800 text-emerald-500 focus:ring-0 focus:ring-offset-0 bg-zinc-900 mt-0.5"
                      />
                      <div>
                        <div className="text-xs font-bold text-white">Identity Credentials Authenticated</div>
                        <p className="text-[10px] text-zinc-400 mt-0.5">Verification matches landlord deed registration scan.</p>
                      </div>
                    </button>

                    {/* Check 2 */}
                    <button
                      onClick={() => handleToggleCheck('escrow')}
                      className={`w-full p-3 rounded-xl border text-left flex items-start gap-3 transition-all ${
                        checklist.escrow 
                          ? 'bg-emerald-500/5 border-emerald-500/30 text-emerald-400' 
                          : 'bg-zinc-950/40 border-zinc-800 hover:border-zinc-700/80 text-zinc-300'
                      }`}
                    >
                      <input 
                        type="checkbox" 
                        checked={checklist.escrow}
                        readOnly
                        className="rounded border-zinc-800 text-emerald-500 focus:ring-0 focus:ring-offset-0 bg-zinc-900 mt-0.5"
                      />
                      <div>
                        <div className="text-xs font-bold text-white">Escrow Voucher Balance Locked</div>
                        <p className="text-[10px] text-zinc-400 mt-0.5">Verify that the collateral matches the specified amount.</p>
                      </div>
                    </button>

                    {/* Check 3 */}
                    <button
                      onClick={() => handleToggleCheck('permit')}
                      className={`w-full p-3 rounded-xl border text-left flex items-start gap-3 transition-all ${
                        checklist.permit 
                          ? 'bg-emerald-500/5 border-emerald-500/30 text-emerald-400' 
                          : 'bg-zinc-950/40 border-zinc-800 hover:border-zinc-700/80 text-zinc-300'
                      }`}
                    >
                      <input 
                        type="checkbox" 
                        checked={checklist.permit}
                        readOnly
                        className="rounded border-zinc-800 text-emerald-500 focus:ring-0 focus:ring-offset-0 bg-zinc-900 mt-0.5"
                      />
                      <div>
                        <div className="text-xs font-bold text-white">Business Permit & Tax Valid</div>
                        <p className="text-[10px] text-zinc-400 mt-0.5">Verify Lessor is fully registered with trade bureau.</p>
                      </div>
                    </button>

                    {/* Check 4 */}
                    <button
                      onClick={() => handleToggleCheck('aml')}
                      className={`w-full p-3 rounded-xl border text-left flex items-start gap-3 transition-all ${
                        checklist.aml 
                          ? 'bg-emerald-500/5 border-emerald-500/30 text-emerald-400' 
                          : 'bg-zinc-950/40 border-zinc-800 hover:border-zinc-700/80 text-zinc-300'
                      }`}
                    >
                      <input 
                        type="checkbox" 
                        checked={checklist.aml}
                        readOnly
                        className="rounded border-zinc-800 text-emerald-500 focus:ring-0 focus:ring-offset-0 bg-zinc-900 mt-0.5"
                      />
                      <div>
                        <div className="text-xs font-bold text-white">Compliance & AML Clearance</div>
                        <p className="text-[10px] text-zinc-400 mt-0.5">Anti-Money Laundering checks checked and validated.</p>
                      </div>
                    </button>
                  </div>
                </div>
              )}

              {/* Tab 2: Comments & Annotations */}
              {activeTab === 'comments' && (
                <div className="space-y-4 flex flex-col h-full">
                  <div className="space-y-1">
                    <h4 className="text-xs font-black text-white flex items-center gap-1.5 uppercase tracking-wide">
                      <MessageSquare size={15} className="text-orange-500" /> Compliance Audit Trail
                    </h4>
                    <p className="text-[11px] text-zinc-400">Post remarks and annotations relative to the verification details.</p>
                  </div>

                  {/* Comments Feed */}
                  <div className="space-y-3 max-h-64 overflow-y-auto pr-1">
                    {comments.map((cmt) => (
                      <div key={cmt.id} className="p-3 bg-zinc-950 border border-zinc-800 rounded-xl space-y-1">
                        <div className="flex items-center justify-between text-[10px]">
                          <span className="font-extrabold text-orange-400">{cmt.author}</span>
                          <span className="text-zinc-500">{cmt.time}</span>
                        </div>
                        <div className="text-[9px] font-mono text-zinc-400 uppercase tracking-widest leading-none mt-0.5">
                          {cmt.role}
                        </div>
                        <p className="text-[11px] text-zinc-300 mt-1.5 leading-relaxed bg-zinc-900/50 p-2 rounded-lg border border-zinc-850">
                          {cmt.text}
                        </p>
                      </div>
                    ))}
                  </div>

                  {/* Add Remark Form */}
                  <form onSubmit={handleAddComment} className="pt-2 border-t border-zinc-800/80 space-y-2">
                    <textarea
                      placeholder="Type a compliance note or annotation..."
                      value={newComment}
                      onChange={(e) => setNewComment(e.target.value)}
                      rows={3}
                      className="w-full bg-zinc-950 border border-zinc-850 text-white rounded-xl p-3 text-xs focus:outline-none focus:border-orange-500 placeholder-zinc-500 focus:ring-0"
                    />
                    <button
                      type="submit"
                      className="w-full py-2 bg-zinc-800 hover:bg-zinc-700 text-white border border-zinc-700 hover:border-zinc-600 rounded-xl text-xs font-bold transition-colors"
                    >
                      Post Annotation Note
                    </button>
                  </form>
                </div>
              )}

              {/* Tab 3: Cryptographic Ledger Registry */}
              {activeTab === 'crypto' && (
                <div className="space-y-5">
                  <div className="space-y-1">
                    <h4 className="text-xs font-black text-white flex items-center gap-1.5 uppercase tracking-wide">
                      <Fingerprint size={15} className="text-orange-500" /> Digital Ledger Telemetry
                    </h4>
                    <p className="text-[11px] text-zinc-400">Cryptographical parameters verification and validation stamps.</p>
                  </div>

                  <div className="p-4 bg-zinc-950 rounded-2xl border border-zinc-800 space-y-4 text-xs">
                    
                    {/* Hash Code */}
                    <div className="space-y-1">
                      <div className="flex items-center justify-between text-[10px] font-mono text-zinc-400 uppercase font-bold">
                        <span>SHA-256 Ledger Fingerprint</span>
                        <button 
                          onClick={handleCopyHash}
                          className="hover:text-white flex items-center gap-1 transition-colors"
                        >
                          {copiedHash ? <Check size={11} className="text-emerald-400" /> : <Copy size={11} />}
                          {copiedHash ? 'Copied' : 'Copy'}
                        </button>
                      </div>
                      <div className="p-2.5 bg-zinc-900 border border-zinc-850 rounded-xl font-mono text-[9px] break-all text-zinc-300 select-all leading-relaxed">
                        sha256-8F9083BA571E9023AC{selectedRow.id}E47F20E90312FDC8B95A2
                      </div>
                    </div>

                    {/* Metadata items */}
                    <div className="space-y-2 border-t border-zinc-800/80 pt-3">
                      <div className="flex justify-between items-center text-[10px]">
                        <span className="text-zinc-500 font-bold uppercase tracking-wider">Signing Authority</span>
                        <span className="text-zinc-300 font-mono">SellBuy Root CA v2</span>
                      </div>
                      
                      <div className="flex justify-between items-center text-[10px]">
                        <span className="text-zinc-500 font-bold uppercase tracking-wider">Public Key Fingerprint</span>
                        <span className="text-zinc-300 font-mono">ED25519-90:AF:71:02</span>
                      </div>

                      <div className="flex justify-between items-center text-[10px]">
                        <span className="text-zinc-500 font-bold uppercase tracking-wider">Verification Epoch</span>
                        <span className="text-zinc-300 font-mono">1784391204 Epoch</span>
                      </div>

                      <div className="flex justify-between items-center text-[10px]">
                        <span className="text-zinc-500 font-bold uppercase tracking-wider">Escrow Ledger Link</span>
                        <span className="text-emerald-400 font-mono font-bold flex items-center gap-1">
                          <CheckCircle2 size={10} /> Active Node Synchronized
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="p-4.5 bg-gradient-to-r from-orange-600/10 to-amber-600/10 border border-orange-500/20 rounded-2xl text-[11px] text-zinc-300 leading-relaxed flex items-start gap-3">
                    <ShieldAlert size={18} className="text-orange-400 mt-0.5 shrink-0" />
                    <div>
                      <span className="font-extrabold text-white">Ledger Tamper Prevention:</span> Any custom modifications to this agreement will invalidate the cryptographical token signature instantly, notifying the Compliance firewall.
                    </div>
                  </div>
                </div>
              )}

            </div>

            {/* Verification Action Drawer (Footer of right-pane) */}
            <div className="p-4 bg-zinc-950 border-t border-zinc-800 space-y-2.5">
              <div className="text-[10px] text-zinc-400 uppercase tracking-widest font-extrabold font-mono text-center">
                Review Decision Center
              </div>
              
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => handleAuditAction('Rejected')}
                  className="py-2.5 px-3 bg-red-600/10 hover:bg-red-600/20 text-red-400 border border-red-500/20 hover:border-red-500/40 rounded-xl text-xs font-black transition-all flex items-center justify-center gap-1.5"
                >
                  <AlertTriangle size={14} /> Flag / Decline
                </button>

                <button
                  type="button"
                  onClick={() => handleAuditAction('Approved')}
                  className={`py-2.5 px-3 rounded-xl text-xs font-black transition-all flex items-center justify-center gap-1.5 ${
                    progressPercent === 100
                      ? 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-lg shadow-emerald-600/20'
                      : 'bg-zinc-800 text-zinc-500 border border-zinc-700/50 cursor-not-allowed'
                  }`}
                  title={progressPercent === 100 ? 'Approve & Release escrow listing' : 'Verify all compliance checklist items to enable'}
                >
                  <ShieldCheck size={14} /> Certify & Approve
                </button>
              </div>

              <div className="text-[10px] text-center text-zinc-500 leading-tight">
                Reviewing updates will write back a sealed audit ledger block to the database.
              </div>
            </div>

          </div>

        </div>

        {/* Modal Bottom Status bar */}
        <div className="p-3 bg-zinc-950 border-t border-zinc-900 flex items-center justify-between text-[11px] text-zinc-400">
          <div className="flex items-center gap-2 font-mono">
            <CheckCircle2 size={13} className="text-emerald-400" />
            <span>Digital Ledger Block validated. Checksum matches database.</span>
          </div>

          <button 
            onClick={onClose} 
            className="px-4 py-1.5 bg-zinc-900 hover:bg-zinc-800 text-zinc-300 rounded-lg font-bold border border-zinc-800 transition-colors"
          >
            Close Audit
          </button>
        </div>
      </div>
    </div>
  );
}
