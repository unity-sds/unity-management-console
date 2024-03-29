import type {MarketplaceMetadata} from "./unity-cs-manager/protobuf/marketplace";

export interface NodeGroupType {
	name: string;
	settings: {
		MinNodes: number;
		MaxNodes: number;
		DesiredNodes: number;
		InstanceType: string;
	};
}


export class OrderLine {
	constructor(public product: MarketplaceMetadata, public quantity: number) {}

	get total(): number {
		return 0;
	}
}

export class Application {
	constructor(public app: MarketplaceMetadata) {}
}

export class Order {
	private lines = new Map<number, OrderLine>();

	constructor(initialLines?: OrderLine[]) {
		//if (initialLines) initialLines.forEach((ol) => this.lines.set(ol.product.Id, ol));
	}

	public addProduct(prod: MarketplaceMetadata, quantity: number) {
		// if (this.lines.has(prod.Id)) {
		// 	if (quantity === 0) {
		// 		this.removeProduct(prod.Id);
		// 	} else {
		// 		const orderLine = this.lines.get(prod.Id);
		// 		if (orderLine)
		// 			// map.get() may return undefined
		// 			orderLine.quantity += quantity;
		// 	}
		// } else {
		// 	this.lines.set(prod.Id, new OrderLine(prod, quantity));
		// }
	}

	public removeProduct(id: number) {
		this.lines.delete(id);
	}

	get orderLines(): OrderLine[] {
		return [...this.lines.values()];
	}

	get productCount(): number {
		return [...this.lines.values()].reduce((total, ol) => (total += ol.quantity), 0);
	}

	get total(): number {
		return [...this.lines.values()].reduce((total, ol) => (total += ol.total), 0);
	}
}
