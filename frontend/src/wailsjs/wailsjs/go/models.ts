export namespace handler {
	
	export class AdjustStockRequest {
	    product_id: number;
	    quantity: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new AdjustStockRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.quantity = source["quantity"];
	        this.notes = source["notes"];
	    }
	}
	export class CategoryResponse {
	    id: number;
	    name: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class CreateCategoryRequest {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class CreateProductRequest {
	    category_id?: number;
	    name: string;
	    code: string;
	    barcode: string;
	    price: number;
	    stock: number;
	    min_stock: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateProductRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.name = source["name"];
	        this.code = source["code"];
	        this.barcode = source["barcode"];
	        this.price = source["price"];
	        this.stock = source["stock"];
	        this.min_stock = source["min_stock"];
	    }
	}
	export class SaleItemRequest {
	    product_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new SaleItemRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class CreateSaleRequest {
	    items: SaleItemRequest[];
	
	    static createFrom(source: any = {}) {
	        return new CreateSaleRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], SaleItemRequest);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProductResponse {
	    id: number;
	    category_id?: number;
	    name: string;
	    code: string;
	    barcode: string;
	    price: number;
	    stock: number;
	    min_stock: number;
	    active: boolean;
	    low_stock: boolean;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category_id = source["category_id"];
	        this.name = source["name"];
	        this.code = source["code"];
	        this.barcode = source["barcode"];
	        this.price = source["price"];
	        this.stock = source["stock"];
	        this.min_stock = source["min_stock"];
	        this.active = source["active"];
	        this.low_stock = source["low_stock"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	
	export class SaleItemResponse {
	    id: number;
	    product_id: number;
	    quantity: number;
	    unit_price: number;
	    subtotal: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new SaleItemResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_id = source["product_id"];
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.subtotal = source["subtotal"];
	        this.created_at = source["created_at"];
	    }
	}
	export class SaleResponse {
	    id: number;
	    total: number;
	    items: SaleItemResponse[];
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new SaleResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.total = source["total"];
	        this.items = this.convertValues(source["items"], SaleItemResponse);
	        this.created_at = source["created_at"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StockMovementResponse {
	    id: number;
	    product_id: number;
	    quantity: number;
	    movement_type: string;
	    notes: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new StockMovementResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_id = source["product_id"];
	        this.quantity = source["quantity"];
	        this.movement_type = source["movement_type"];
	        this.notes = source["notes"];
	        this.created_at = source["created_at"];
	    }
	}
	export class UpdateCategoryRequest {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class UpdateProductRequest {
	    category_id?: number;
	    name: string;
	    code: string;
	    barcode: string;
	    price: number;
	    min_stock: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProductRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.name = source["name"];
	        this.code = source["code"];
	        this.barcode = source["barcode"];
	        this.price = source["price"];
	        this.min_stock = source["min_stock"];
	    }
	}

}

