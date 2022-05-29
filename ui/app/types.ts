export type ApiResponse<T> = {
    success: boolean;
    data: T;
};

export type ProductTypeT = 'food' | 'drink' | 'pastry';

export type ProductT = {
    id: number;
    name: string;
    price: number;
    type: ProductTypeT;
    discontinued: boolean;
    sold_out: boolean;
};

export type ProductAggregateT = {
    product: ProductT;
    amount: number;
}

export type ProductMapT = {
    [key: string]: ProductAggregateT;
};

export type OrderT = {
    id: number;
    cancelled: boolean;
    created_at: string;
    products: ProductT[];
};

export type RichOrderT = {
    order_id: number;
    total: number;
    order: OrderT;
};

export type StationT = {
    id: number;
    name: string;
    products: ProductT[];
    created_at: Date;
    updated_at: Date;
};
