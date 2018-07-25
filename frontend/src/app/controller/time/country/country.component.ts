import { Component, OnInit, ViewChild } from '@angular/core';
import { Comparator, State, SortOrder}  from "clarity-angular";

import { PageSet }              from '../../../model/base/basic';
import { Country }              from '../../../model/time/country';
import { CountryService }       from '../../../service/time/country.service';
import { CountrySaveComponent } from './save/save.component';

import { CustomComparator }                             from '../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../shared/utils';
import { PageSize }                                     from '../../../shared/const';

@Component({
  selector: 'time-country',
  templateUrl: './country.component.html',
  styleUrls: ['./country.component.css']
})
export class CountryComponent implements OnInit {
  @ViewChild(CountrySaveComponent)
  saveCountry: CountrySaveComponent;

  countries: Country[];
  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private countryService: CountryService,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: any): void {
    let country = new Country();
    country = loadPageFilterSort<Country>(country, state);
    this.pageSet.Current = country.Page.Current;
    this.currentState = state;
    this.countryService.Count(country).subscribe(count => {
      this.pageSet.Count = count;
      this.countryService.List(country).subscribe(res => {
        this.countries = res;
      })
    })
  }

  saved(savedCountry: Country): void {
    if (savedCountry.Id) {
      this.load(this.currentState);
    } else {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveCountry.New(id);
  }

  delete(country: Country): void {
    this.countryService.Delete(country.Id).subscribe(res => {
      let state = deleteState(this.pageSet, this.currentState, 1);
      this.load(state);
    })
  }
}
