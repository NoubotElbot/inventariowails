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

