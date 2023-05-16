export interface Product {
	  Id: number;
	  Name: string;
	  Description: string;
	  Category: string;
	  Tags: string[];
    IamRoles: any[];
    Package: string;
    ManagedDependencies: ManagedDependency[];
    Backend: string;
    DefaultDeployment: DefaultDeployment;
    Version: string
    Branch: string
}

interface ManagedDependency {
    Eks: {
        MinimumVersion: string;
    };
}

interface DefaultDeployment {
    Variables: {
        [key: string]: string;
    };
    EksSpec: {
        NodeGroups: NodeGroup[];
    };
}

interface NodeGroup {
    [key: string]: {
        MinNodes: number;
        MaxNodes: number;
        DesiredNodes: number;
        InstanceType: string;
    };
}

export class OrderLine {
	constructor(public product: Product, public quantity: number) {}

	get total(): number {
		  //return this.product.price * this.quantity;
      return 0;
	}
}

export class Application {
    constructor(public app: Product) {}

}

export interface AppInstall {
    
}

export interface Extensions {
    
}

export class Install {
    constructor(public install?: AppInstall[], public extenstions?: Extensions) {}
}
export class Order {
	private lines = new Map<number, OrderLine>();

	constructor(initialLines?: OrderLine[]) {
		if (initialLines) initialLines.forEach((ol) => this.lines.set(ol.product.Id, ol));
	}

	public addProduct(prod: Product, quantity: number) {
		if (this.lines.has(prod.Id)) {
			if (quantity === 0) {
				this.removeProduct(prod.Id);
			} else {
				const orderLine = this.lines.get(prod.Id);
				if (orderLine)
					// map.get() may return undefined
					orderLine.quantity += quantity;
			}
		} else {
			this.lines.set(prod.Id, new OrderLine(prod, quantity));
		}
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
