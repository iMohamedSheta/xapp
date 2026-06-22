import type { LucideIcon } from 'lucide-vue-next';
// import type { Config } from 'ziggy-js';
import { toast, ToastT } from 'vue-sonner';
import { DateValue } from 'reka-ui';
import { SweetAlertCustomClass } from './../../../node_modules/sweetalert2/sweetalert2.d';
import { SweetAlertOptions } from 'sweetalert2';

import type { PermissionKey } from './permissions';

export interface Auth {
  user: User;
  is_impersonating: boolean;
  permissions?: Partial<Record<PermissionKey, boolean>>;
}

export interface Settings {
  active_workspace?: string;
  user_pages?: PageModel[];
}

export interface BreadcrumbItem {
  title: string;
  href: string;
}

export interface NavItem {
  activeKey?: string;
  title: string;
  href: string;
  icon?: LucideIcon;
  isActive?: boolean;
  permission?: PermissionKey;
  children?: NavItem[];
  role?: Role;
}

export interface ToastMessage {
  title: string
  description?: string
  action?: {
    label: string
  }
}

export interface LottieProps {
  autoplay?: boolean
  loop?: boolean
  speed?: number
  playOnHover?: boolean
  backgroundColor?: string
  animationId?: string
  layout?: any
  segment?: [number, number]
  width?: string
  height?: string
}

export interface LottieFlash {
  src: string
  timer?: number
  class: string
  lottieProps?: LottieProps
  closeOnClick?: boolean
  closeOnEscape?: boolean
  autoClose?: boolean
}

export interface Notification {
  id: number
  title: string
  message: string
  time: string
  link: string
  read: boolean
  open: boolean
  type: string
}

export interface FlashMessage {
  toast?: ToastT[];
  sweetalert?: SweetAlertOptions[];
  lottie?: LottieFlash;
  drawer?: string;
  redirect_url?: string;
}

export type AppPageProps<T extends Record<string, unknown> = Record<string, unknown>> = T & {
  name: string;
  quote: { message: string; author: string };
  auth: Auth;
  notification: { notifications: Notification[], count: number };
  settings: Settings
  flash?: FlashMessage;
  debug?: Debug;
  sidebarOpen: boolean;
};

export type Role = 'super_manager' | 'manager' | 'super_admin' | 'admin' | 'distributor' | 'client'

export type BalanceSystem = 'until_finish' | 'monthly_auto'
export type ProfitMethod = 'fixed' | 'percentage'


export interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  avatar?: string;
  email_verified_at: string | null;
  agency_id: number;
  role: Role;
  status: number;
  deleted_at: string | null;
  created_at: string;
  updated_at: string;

  permissions?: string[];
}

export interface Order {
  id: number;
  agency_id: number;
  workspace_id: number;
  user_id?: number | null;
  orderable_type: string;
  orderable_id: number;
  quantity: number;
  unit_price: number;
  total_price: number;
  currency: string;
  status: number;
  created_at: string;
  updated_at: string;
  user?: User;
  transaction?: Transaction;
  invoice?: Invoice;
  event?: Event;
}



export interface PlanPrice {
  id?: number;
  plan_id?: number;
  price: number;
  currency: string;
  created_at?: string;
  updated_at?: string;
}

export interface PlanSettings {
  id?: number;
  plan_id?: number;
  expire_action: 'block' | 'downgrade'
  downgrade_to_plan: number | null
  grace_period_days: number
  created_at?: string;
  updated_at?: string;

  view_price?: number;
}

export interface Plan {
  id: number;
  name: string;
  features: string[];
  is_active: boolean;
  popular: boolean;
  created_at: string;
  updated_at: string;
  plan_setting: PlanSettings;
  plan_prices: PlanPrice[];
}

export type SubscriptionStatus = 'active' | 'expired' | 'canceled' | 'trial';

export type BillingCycle = 'monthly' | 'yearly' | 'onetime';

export interface Subscription {
  id: number;
  agency_id: number;
  plan_id: number;
  start_date: string; // ISO date string
  end_date: string;
  status: SubscriptionStatus;
  price: number;
  original_price: number;
  currency: string;
  auto_renew: boolean;
  billing_cycle: BillingCycle;
  expire_action: string;
  downgrade_to_plan?: number | null;
  grace_period_days: number;
  created_at: string;
  updated_at: string;
  plan?: Plan;
}


export type InvoiceableType = "subscription";

export type InvoiceType =
  | "subscription"
  | "upgrade"
  | "downgrade"
  | "renew_subscription"
  | "manual_charge"
  | "refund"
  | "manual_debit"
  | "manual_charge";

export type InvoiceStatus =
  | "pending"
  | "paid"
  | "partially_paid"
  | "failed"
  | "refunded"
  | "canceled";

export interface Invoice {
  id: number;
  agency_id: number;
  workspace_id?: number | null;
  user_id?: number | null;
  order_id?: number | null;
  transaction_id?: number | null;

  invoiceable_type: InvoiceableType;
  invoiceable_id: number;

  invoice_number: string;
  type: InvoiceType;
  status: InvoiceStatus;

  amount: number;
  paid: number;
  currency: string;

  due_date?: string | null; // ISO string
  paid_at?: string | null;  // ISO string

  notes?: string | null;
  metadata?: any; // JSON object

  created_at: string; // ISO string
  updated_at: string; // ISO string

  // Relationships (optional)
  user?: User;
  order?: Order;
  transaction?: Transaction;

  customer_fullname?: string;
  agency_name?: string;
}

// models.ts
export interface Transaction {
  id: number;
  agencyId: number;
  workspaceId: number;
  userId?: number | null;
  orderId: number;

  // Gateway info
  gateway: string;
  gatewayTxId?: string;
  gatewayOrderId?: string;
  merchantOrderId?: string;

  // Amounts
  amountCents: number;
  currency: string;
  capturedAmount: number;
  refundedAmount: number;

  // Status
  success: boolean;
  status?: string;
  responseCode?: string;
  message?: string;

  // Payment Method Info (non-sensitive)
  paymentMethod?: string;
  cardType?: string;
  cardLast4?: string;

  // Gateway Metadata
  authorizeId?: string;
  receiptNo?: string;
  batchNo?: string;
  gatewayRefId?: string;

  // Customer Snapshot
  customerEmail?: string;
  customerPhone?: string;
  customerName?: string;

  // Raw JSON response
  rawResponse?: any; // JSON object

  // Timestamps
  createdAt: string; // ISO string
  updatedAt: string; // ISO string

  // Relations
  order?: Order;
}


export interface PaginationMeta {
  total_count: number;
  last_page: number;
  next_page: number;
  prev_page: number;
  per_page: number;
  current_page: number;
}

export interface PageTemplate {
  id: number;
  title: string;
  image_path: string;
  path: string;
  description: string;
}

export interface Debug {
  queries: Array<{
    buildDuration: string;
    execDuration: string;
    sql: string;
    rawSql: string;
    bindings: any[];
    sourceFile: string;
    sourceLine: number;
    expanded?: boolean;
  }>;
  request: {
    id: string;
    host: string;
    method: string;
    query: string;
    path: string;
    action_duration: string;
    protocol: string;
    referer: string;
    user_agent: string;
    client_ip: string;
    status: number;
    memory_usage: string;
    memory_usage_for_request: string;
    user: any | null;
  };
  route: {
    action: string;
    handler: string;
    method: string;
    middleware: string[];
    short: string;
    uri: string;
  };
}


export type BreadcrumbItemType = BreadcrumbItem;
