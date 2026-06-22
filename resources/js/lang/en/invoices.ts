export default {
    title: "My Invoices",
    breadcrumbs: "Invoices",
    stats: {
        total: "Total",
        total_desc: "Total Invoices",
        paid: "Paid",
        paid_desc: "Paid Amount",
        unpaid: "Unpaid",
        unpaid_desc: "Unpaid Amount",
        unpaid_count_title: "Unpaid Invoices",
        unpaid_count_desc: "Unpaid Invoice Count",
    },
    filters: {
        customer_type: "Customer Type",
        select_customer_type: "Select Customer Type",
        type: "Type",
        select_type: "Select Type",
        status: "Status",
        select_status: "Select Status",
        from_date: "From Date",
        to_date: "To Date",
        clear_all: "Clear All Filters",
    },
    table: {
        invoice_number: "Invoice #",
        customer_name: "Customer Name",
        type: "Type",
        amount: "Total Amount",
        paid: "Paid",
        remaining: "Remaining",
        status: "Status",
        paid_at: "Paid At",
        created_at: "Created At",
    },
    actions: {
        pay: "Pay Invoice",
        print: "Print Invoice",
    },
    empty: {
        title: "No Invoices Found",
        description: "No invoices matches your search criteria. Try changing filters.",
    }
}
