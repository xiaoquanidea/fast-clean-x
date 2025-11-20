export namespace models {
	
	export class ScanRule {
	    name: string;
	    description: string;
	    targetDirs: string[];
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScanRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.targetDirs = source["targetDirs"];
	        this.enabled = source["enabled"];
	    }
	}
	export class Config {
	    scanPaths: string[];
	    ignorePatterns: string[];
	    scanRules: ScanRule[];
	    // Go type: time
	    lastScanTime: any;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanPaths = source["scanPaths"];
	        this.ignorePatterns = source["ignorePatterns"];
	        this.scanRules = this.convertValues(source["scanRules"], ScanRule);
	        this.lastScanTime = this.convertValues(source["lastScanTime"], null);
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
	export class ScanItem {
	    path: string;
	    projectPath: string;
	    projectName: string;
	    type: string;
	    size: number;
	    sizeReadable: string;
	    fileCount: number;
	    // Go type: time
	    lastModified: any;
	    selected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScanItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.projectPath = source["projectPath"];
	        this.projectName = source["projectName"];
	        this.type = source["type"];
	        this.size = source["size"];
	        this.sizeReadable = source["sizeReadable"];
	        this.fileCount = source["fileCount"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.selected = source["selected"];
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
	export class ScanResult {
	    items: ScanItem[];
	    totalSize: number;
	    totalCount: number;
	    // Go type: time
	    scanTime: any;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], ScanItem);
	        this.totalSize = source["totalSize"];
	        this.totalCount = source["totalCount"];
	        this.scanTime = this.convertValues(source["scanTime"], null);
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

}

