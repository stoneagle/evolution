import { Component, OnInit, Output, Input, EventEmitter, ViewChild } from '@angular/core';
import { NgForm  }                                                   from '@angular/forms';

import { AssetSource }    from '../../../../model/quant/config';
import { ConfigService  } from '../../../../service/quant/config.service';

@Component({
  selector: 'quant-config-asset-source',
  templateUrl: './asset-source.component.html',
  styleUrls: ['./asset-source.component.css']
})

export class ConfigAssetSourceComponent implements OnInit {
  assetSource: AssetSource = new AssetSource;

  selectMap: Map<string, Map<string, Map<string, string[]>>> = new Map();
  mainAndSubMap: Map<string, string[]> = new Map();

  assetMap: Map<string, string> = new Map();
  typeMap: Map<string, string> = new Map();
  mainMap: Map<string, string> = new Map();
  subMap: Map<string, string> = new Map();

  @Output() sync = new EventEmitter<any>();
  @Input() modelOpened: boolean;

  @ViewChild('syncForm')
  currentForm: NgForm;

  constructor(
    private configService: ConfigService,
  ) { }

  New(): void {
  }            

  Cancel(): void {
    this.sync.emit(false);
  }

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

  AssetOnChange(asset: string) {
    let typeMap = this.selectMap.get(asset);
    this.typeMap = new Map();
    typeMap.forEach((mainMap, ctype) => {
      this.typeMap.set(ctype, ctype);
    });
    this.assetSource.AssetString = asset;
    this.assetSource.Type = this.typeMap.keys().next().value;
    this.TypeOnChange(this.typeMap.keys().next().value);
  }

  TypeOnChange(ctype) {
    let mainMap = this.selectMap.get(this.assetSource.AssetString).get(ctype);
    this.mainMap = new Map();
    mainMap.forEach((subMap, main) => {
      this.mainMap.set(main, main);
    });
    this.assetSource.Type = ctype; 
    this.assetSource.Main = this.mainMap.keys().next().value;
    this.MainOnChange(this.mainMap.keys().next().value);
  }

  MainOnChange(main) {
    let subMap = this.selectMap.get(this.assetSource.AssetString).get(this.assetSource.Type).get(main);
    this.subMap = new Map();
    subMap.forEach((sub, key) => {
      this.subMap.set(sub, sub);
    });
    this.assetSource.Sub = this.subMap.keys().next().value;
  }

  getKeys(map) {
    return Array.from(map.keys());
  }

  Submit(): void {
    this.sync.emit(this.assetSource);
  }
}
