import { Component, OnInit, Output, EventEmitter } from '@angular/core'; 
import { Item } from '../../../model/business/item';
import { ItemService  } from '../../../service/business/item.service';
import { ConfigService } from '../../../service/config/config.service';

@Component({
  selector: 'sync-item',
  templateUrl: './sync.component.html',
  styleUrls: ['./sync.component.css']
})

export class SyncItemComponent implements OnInit {
  item: Item = new Item;
  modelOpened: boolean = false;
  assetMap: Map<string, string> = new Map();
  selectMap: Map<string, Map<string, Map<string, string[]>>> = new Map();

  typeMap: Map<string, string> = new Map();
  sourceMap: Map<string, string> = new Map();
  subMap: Map<string, string> = new Map();

  sourceAndSubMap: Map<string, string[]> = new Map();


  @Output() sync = new EventEmitter<boolean>();

  constructor(
    private itemService: ItemService,
    private configService: ConfigService,
  ) { }

  ngOnInit() {
    this.assetMap = new Map();
    this.configService.AssetList()
      .subscribe(res => {
        this.assetMap = res;
        this.assetMap.forEach((key, asset) => {
          this.configService.TypeList(asset)
            .subscribe(res => {
              this.selectMap.set(asset, res);
            })
        })
      })
  }

  // AssetOnChange(asset: string) {
  //   let typeMap = this.selectMap.get(asset);
  //   this.typeMap = new Map();
  //   typeMap.forEach((sourceMap, ctype) => {
  //     this.typeMap.set(ctype, ctype);
  //   });
  //   this.item.AssetString = asset;
  //   this.item.Type = this.typeMap.keys().next().value;
  //   this.TypeOnChange(this.typeMap.keys().next().value);
  // }

  // TypeOnChange(ctype) {
  //   let sourceMap = this.selectMap.get(this.item.AssetString).get(ctype);
  //   this.sourceMap = new Map();
  //   sourceMap.forEach((subMap, source) => {
  //     this.sourceMap.set(source, source);
  //   });
  //   this.item.Type = ctype; 
  //   this.item.Source = this.sourceMap.keys().next().value;
  //   this.SourceOnChange(this.sourceMap.keys().next().value);
  // }

  // SourceOnChange(source) {
  //   let subMap = this.selectMap.get(this.item.AssetString).get(this.item.Type).get(source);
  //   this.subMap = new Map();
  //   subMap.forEach((sub, key) => {
  //     this.subMap.set(sub, sub);
  //   });
  //   this.item.Sub = this.subMap.keys().next().value;
  // }

  New(id?: number): void {
    this.item = new Item();
    this.modelOpened = true;
  }            

  Sync(): void {
    this.itemService.Sync(this.item)
    .subscribe(res => {
      this.modelOpened = false;
      if (res) {
        this.sync.emit(true);
      } else {
        this.sync.emit(false);
      }
    })
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
