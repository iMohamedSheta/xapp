export default {
    title: "فواتيري",
    breadcrumbs: "الفواتير",
    stats: {
        total: "الإجمالي",
        total_desc: "إجمالي الفواتير",
        paid: "مدفوع",
        paid_desc: "المبلغ المدفوع",
        unpaid: "غير مدفوع",
        unpaid_desc: "المبلغ الغير مدفوع",
        unpaid_count_title: "الفواتير الغير مدفوعة",
        unpaid_count_desc: "فاتورة غير مدفوعة",
    },
    filters: {
        customer_type: "نوع العميل",
        select_customer_type: "اختر نوع العميل",
        type: "النوع",
        select_type: "اختر النوع",
        status: "الحالة",
        select_status: "اختر الحالة",
        from_date: "من تاريخ",
        to_date: "إلى تاريخ",
        clear_all: "مسح جميع الفلاتر",
    },
    table: {
        invoice_number: "رقم الفاتورة",
        customer_name: "اسم العميل",
        type: "النوع",
        amount: "المبلغ الإجمالي",
        paid: "المدفوع",
        remaining: "المتبقي",
        status: "الحالة",
        paid_at: "تاريخ الدفع",
        created_at: "تاريخ الإنشاء",
    },
    actions: {
        pay: "دفع الفاتورة",
        print: "طباعة الفاتورة",
    },
    empty: {
        title: "لا توجد فواتير",
        description: "لم يتم العثور على أي فواتير. جرب تغيير الفلاتر أو البحث عن شيء آخر.",
    }
}
