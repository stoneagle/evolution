import { Component, OnInit, ViewChild } from '@angular/core';

import { Classify }             from '../../../model/quant/classify';
import { ClassifyService }      from '../../../service/quant/classify.service';
import { ItemService }          from '../../../service/quant/item.service';
import { AssetSource }          from '../../../model/quant/config';
import { ConfigAssetSourceComponent } from '../config/asset-source/asset-source.component';

@Component({
  selector: 'quant-classify',
  templateUrl: './classify.component.html',
  styleUrls: ['./classify.component.css']
})
export class ClassifyComponent implements OnInit {
  @ViewChild(ConfigAssetSourceComponent)
  assetSourceComponent: ConfigAssetSourceComponent;

  classifys: Classify[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  syncModelOpened: boolean = false;

  constructor(
    private classifyService: ClassifyService,
    private itemService: ItemService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  Sync($event: any): void {
    if ($event) {
      this.classifyService.Sync($event).subscribe(res => {
        this.syncModelOpened = false;
        this.refresh();
      })
    } else {
      this.syncModelOpened = false;
    }
  }

  syncItem(classify: Classify): void {
    this.itemService.SyncClassify(classify).subscribe(res => {
    })
  }

  openSyncModel(): void {
    this.syncModelOpened = true;
    this.assetSourceComponent.New()
  }

  delete(classify: Classify): void {
    this.classifyService.Delete(classify.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    this.classifyService.List(null).subscribe(res => {
      this.totalCount = res.length;
      this.classifys = res.slice(from, to);
    })
  }
}
